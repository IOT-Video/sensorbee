package bql

import (
	"errors"
	"fmt"
	"pfi/sensorbee/sensorbee/bql/execution"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"sync"
	"sync/atomic"
)

type TopologyBuilder struct {
	topology       core.Topology
	Reg            udf.FunctionManager
	UDSFCreators   udf.UDSFCreatorRegistry
	UDSCreators    udf.UDSCreatorRegistry
	SourceCreators SourceCreatorRegistry
	SinkCreators   SinkCreatorRegistry
}

// TODO: Provide AtomicTopologyBuilder which support building multiple nodes
// in an atomic manner (kind of transactionally)

// NewTopologyBuilder creates a new TopologyBuilder which dynamically creates
// nodes from BQL statements. The target Topology can be shared by
// multiple TopologyBuilders.
//
// TopologyBuilder doesn't support atomic topology building. For example,
// when a user wants to add three statement and the second statement fails,
// only the node created from the first statement is registered to the topology
// and it starts to generate tuples. Others won't be registered.
func NewTopologyBuilder(t core.Topology) (*TopologyBuilder, error) {
	udsfs, err := udf.CopyGlobalUDSFCreatorRegistry()
	if err != nil {
		return nil, err
	}

	udss, err := udf.CopyGlobalUDSCreatorRegistry()
	if err != nil {
		return nil, err
	}

	srcs, err := CopyGlobalSourceCreatorRegistry()
	if err != nil {
		return nil, err
	}
	// node_statuses builtin source can only be registered here because it
	// requires a topology.
	if err := srcs.Register("node_statuses", createNodeStatusSourceCreator(t)); err != nil {
		return nil, err
	}

	sinks, err := CopyGlobalSinkCreatorRegistry()
	if err != nil {
		return nil, err
	}

	tb := &TopologyBuilder{
		topology:       t,
		Reg:            udf.CopyGlobalUDFRegistry(t.Context()),
		UDSFCreators:   udsfs,
		UDSCreators:    udss,
		SourceCreators: srcs,
		SinkCreators:   sinks,
	}
	return tb, nil
}

// TODO: if IDs are shared by distributed processes, they should have a process
// ID of each process in it or they should be generated by a central server
// to be globally unique. Currently, this id is only used temporarily and
// doesn't have to be strictly unique nor recoverable (durable).
var (
	topologyBuilderTemporaryID int64
)

func topologyBuilderNextTemporaryID() int64 {
	return atomic.AddInt64(&topologyBuilderTemporaryID, 1)
}

func (tb *TopologyBuilder) Topology() core.Topology {
	return tb.topology
}

// AddStmt add a node created from a statement to the topology. It returns
// a created node. It returns a nil node when the statement is CREATE STATE.
func (tb *TopologyBuilder) AddStmt(stmt interface{}) (core.Node, error) {
	// TODO: Enable StopOnDisconnect properly

	// check the type of statement
	switch stmt := stmt.(type) {
	case parser.CreateSourceStmt:
		// load params into map for faster access
		paramsMap := tb.mkParamsMap(stmt.Params)

		// check if we know whis type of source
		creator, err := tb.SourceCreators.Lookup(string(stmt.Type))
		if err != nil {
			return nil, err
		}

		// if so, try to create such a source
		source, err := creator.CreateSource(tb.topology.Context(), &IOParams{
			TypeName: string(stmt.Type),
			Name:     string(stmt.Name),
		}, paramsMap)
		if err != nil {
			return nil, err
		}
		return tb.topology.AddSource(string(stmt.Name), source, &core.SourceConfig{
			PausedOnStartup: stmt.Paused == parser.Yes,
		})

	case parser.CreateStreamAsSelectStmt:
		return tb.createStreamAsSelectStmt(&stmt)

	case parser.CreateSinkStmt:
		// load params into map for faster access
		paramsMap := tb.mkParamsMap(stmt.Params)

		// check if we know whis type of sink
		creator, err := tb.SinkCreators.Lookup(string(stmt.Type))
		if err != nil {
			return nil, err
		}

		// if so, try to create such a sink
		sink, err := creator.CreateSink(tb.topology.Context(), &IOParams{
			TypeName: string(stmt.Type),
			Name:     string(stmt.Name),
		}, paramsMap)
		if err != nil {
			return nil, err
		}
		// we insert a sink, but cannot connect it to
		// any streams yet, therefore we have to keep track
		// of the SinkDeclarer
		return tb.topology.AddSink(string(stmt.Name), sink, nil)

	case parser.CreateStateStmt:
		c, err := tb.UDSCreators.Lookup(string(stmt.Type))
		if err != nil {
			return nil, err
		}

		ctx := tb.topology.Context()
		s, err := c.CreateState(ctx, tb.mkParamsMap(stmt.Params))
		if err != nil {
			return nil, err
		}
		if err := ctx.SharedStates.Add(string(stmt.Name), s); err != nil {
			return nil, err
		}
		return nil, nil

	case parser.UpdateStateStmt:
		ctx := tb.topology.Context()
		state, err := ctx.SharedStates.Get(string(stmt.Name))
		if err != nil {
			return nil, err
		}

		u, ok := state.(core.Updater)
		if !ok {
			return nil, fmt.Errorf("%s cannot be updated", string(stmt.Name))
		}
		return nil, u.Update(tb.mkParamsMap(stmt.Params))

	case parser.UpdateSourceStmt:
		src, err := tb.topology.Source(string(stmt.Name))
		if err != nil {
			return nil, err
		}

		u, ok := src.Source().(core.Updater)
		if !ok {
			return nil, fmt.Errorf("%s cannot be updated", string(stmt.Name))
		}
		return nil, u.Update(tb.mkParamsMap(stmt.Params))

	case parser.UpdateSinkStmt:
		sink, err := tb.topology.Sink(string(stmt.Name))
		if err != nil {
			return nil, err
		}

		u, ok := sink.Sink().(core.Updater)
		if !ok {
			return nil, fmt.Errorf("%s cannot be updated", string(stmt.Name))
		}
		return nil, u.Update(tb.mkParamsMap(stmt.Params))

	case parser.DropSourceStmt:
		_, err := tb.topology.Source(string(stmt.Source))
		if err != nil {
			return nil, err
		}

		return nil, tb.topology.Remove(string(stmt.Source))

	case parser.DropStreamStmt:
		_, err := tb.topology.Box(string(stmt.Stream))
		if err != nil {
			return nil, err
		}

		return nil, tb.topology.Remove(string(stmt.Stream))

	case parser.DropSinkStmt:
		_, err := tb.topology.Sink(string(stmt.Sink))
		if err != nil {
			return nil, err
		}

		return nil, tb.topology.Remove(string(stmt.Sink))

	case parser.DropStateStmt:
		ctx := tb.topology.Context()
		_, err := ctx.SharedStates.Get(string(stmt.State))
		if err != nil {
			return nil, err
		}

		_, err = ctx.SharedStates.Remove(string(stmt.State))
		return nil, err

	case parser.InsertIntoSelectStmt:
		// get the sink to add an input to
		sink, err := tb.topology.Sink(string(stmt.Sink))
		if err != nil {
			return nil, err
		}
		// construct an intermediate box doing the SELECT computation.
		//   INSERT INTO sink SELECT ISTREAM a, b FROM c [RANGE ...] WHERE d
		// becomes
		//   CREATE STREAM (random_string) AS SELECT ISTREAM a, b
		//   FROM c [RANGE ...] WHERE d
		//  + a connection (random_string -> sink)
		tmpName := fmt.Sprintf("sensorbee_tmp_%v", topologyBuilderNextTemporaryID())
		tmpStmt := parser.CreateStreamAsSelectStmt{
			parser.StreamIdentifier(tmpName),
			parser.SelectStmt{
				stmt.EmitterAST,
				stmt.ProjectionsAST,
				stmt.WindowedFromAST,
				stmt.FilterAST,
				stmt.GroupingAST,
				stmt.HavingAST,
			},
		}
		box, err := tb.AddStmt(tmpStmt)
		if err != nil {
			return nil, err
		}
		box.(core.BoxNode).StopOnDisconnect(core.Inbound | core.Outbound)
		box.(core.BoxNode).RemoveOnStop()

		// now connect the sink to that box
		if err := sink.Input(tmpName, nil); err != nil {
			tb.topology.Remove(tmpName)
			return nil, err
		}
		return box, nil

	case parser.InsertIntoFromStmt:
		// get the sink to add an input to
		sink, err := tb.topology.Sink(string(stmt.Sink))
		if err != nil {
			return nil, err
		}
		// now connect the sink to the specified box
		if err := sink.Input(string(stmt.Input), nil); err != nil {
			return nil, err
		}
		return sink, nil

	case parser.PauseSourceStmt:
		src, err := tb.topology.Source(string(stmt.Source))
		if err != nil {
			return nil, err
		}
		if err := src.Pause(); err != nil {
			return nil, err
		}
		return src, nil

	case parser.ResumeSourceStmt:
		src, err := tb.topology.Source(string(stmt.Source))
		if err != nil {
			return nil, err
		}
		if err := src.Resume(); err != nil {
			return nil, err
		}
		return src, nil

	case parser.RewindSourceStmt:
		src, err := tb.topology.Source(string(stmt.Source))
		if err != nil {
			return nil, err
		}
		if err := src.Rewind(); err != nil {
			return nil, err
		}
		return src, nil
	}

	return nil, fmt.Errorf("statement of type %T is unimplemented", stmt)
}

type udsfBox struct {
	f udf.UDSF
}

var (
	_ core.StatefulBox = &udsfBox{}
)

func newUDSFBox(f udf.UDSF) *udsfBox {
	return &udsfBox{
		f: f,
	}
}

func (b *udsfBox) Init(ctx *core.Context) error {
	return nil
}

func (b *udsfBox) Process(ctx *core.Context, t *core.Tuple, w core.Writer) error {
	return b.f.Process(ctx, t, w)
}

func (b *udsfBox) Terminate(ctx *core.Context) error {
	return b.f.Terminate(ctx)
}

func (tb *TopologyBuilder) createStreamAsSelectStmt(stmt *parser.CreateStreamAsSelectStmt) (core.Node, error) {
	// insert a bqlBox that executes the SELECT statement
	outName := string(stmt.Name)
	box := NewBQLBox(&stmt.Select, tb.Reg)
	// add all the referenced relations as named inputs
	dbox, err := tb.topology.AddBox(outName, box, nil)
	if err != nil {
		return nil, err
	}

	removeNodes := true
	var temporaryNodes []string
	defer func() {
		if !removeNodes {
			return
		}

		for _, n := range temporaryNodes {
			tb.topology.Remove(n)
		}
		tb.topology.Remove(outName)
	}()

	connected := map[string]bool{}
	for _, rel := range stmt.Select.Relations {
		switch rel.Type {
		case parser.ActualStream:
			if connected[rel.Name] {
				// this is a self-join (FROM x [RANGE ...] AS a, x [RANGE ...] AS b)
				// and we already have connected x to this box before
				continue
			}
			if err := dbox.Input(rel.Name, &core.BoxInputConfig{
				// For self-join statements like
				//   ... FROM x AS a [RANGE 2 TUPLES],
				//            x AS b [RANGE 3 TUPLES],
				//            y AS c [RANGE 4 TUPLES] ...
				// the bqlBox will have only *two* inputs:
				//  - one is `x` with `InputName: x`,
				//  - one is `y` with `InputName: y`.
				//
				// When a tuple arrives with `InputName: x`, it will be added
				// to the input buffer for `a` and to the input buffer for `b`
				// in the execution plan. When a tuple arrives with `InputName: y`,
				// it will be added to the input buffer for `c`.
				//
				// If we used the rel.Alias here, then we would have to make multiple
				// input connections to the same box, which is not possible.
				InputName: rel.Name,
			}); err != nil {
				return nil, err
			}
			connected[rel.Name] = true

		case parser.UDSFStream:
			// Compute the values of the UDSF parameters (if there was
			// an unusable parameter, as in `udsf(7, col)` this will fail).
			// Note: it doesn't feel exactly right to do this kind of
			// validation here after parsing has been done "successfully",
			// on the other hand the parser should not evaluate expressions
			// (and cannot import the execution package) or make too many
			// semantical checks, so we leave this here for the moment.
			params := make([]data.Value, len(rel.Params))
			for i, expr := range rel.Params {
				p, err := execution.EvaluateFoldable(expr, tb.Reg)
				if err != nil {
					return nil, err
				}
				params[i] = p
			}

			udsfc, err := tb.UDSFCreators.Lookup(rel.Name, len(params))
			if err != nil {
				return nil, err
			}

			decl := udf.NewUDSFDeclarer()
			udsf, err := func() (f udf.UDSF, err error) {
				defer func() {
					if e := recover(); e != nil {
						if er, ok := e.(error); ok {
							err = er
						} else {
							err = fmt.Errorf("cannot create a UDSF: %v", e)
						}
					}
				}()
				return udsfc.CreateUDSF(tb.topology.Context(), decl, params...)
			}()
			if err != nil {
				return nil, err
			}
			if len(decl.ListInputs()) == 0 {
				func() {
					defer func() {
						if e := recover(); e != nil {
							tb.topology.Context().Log().WithField("udsf_name", rel.Name).
								Errorf("Cannot terminate the UDFS due to panic: %v", e)
						}
					}()
					if err := udsf.Terminate(tb.topology.Context()); err != nil {
						tb.topology.Context().ErrLog(err).WithField("udsf_name", rel.Name).
							Error("Cannot terminate the UDSF")
					}
				}()
				return nil, fmt.Errorf("a UDSF '%v' must have at least one input", rel.Name)
			}

			temporaryName := fmt.Sprintf("sensorbee_tmp_udsf_%v", topologyBuilderNextTemporaryID())
			bn, err := tb.topology.AddBox(temporaryName, newUDSFBox(udsf), &core.BoxConfig{
			// TODO: add information of the statement
			})
			if err != nil {
				return nil, err
			}
			temporaryNodes = append(temporaryNodes, temporaryName)

			for input, config := range decl.ListInputs() {
				if err := bn.Input(input, &core.BoxInputConfig{
					InputName: config.InputName,
				}); err != nil {
					return nil, err
				}
			}

			alias := rel.Alias
			if alias == "" {
				alias = rel.Name
			}
			if err := dbox.Input(temporaryName, &core.BoxInputConfig{
				// As opposed to actual streams, for `udsf('s') AS a, udsf('s') AS b`,
				// there will be *multiple* boxes and we will have one connection to
				// udsf 1 (the one aliased to `a`) and udsf 2 (the one aliased to `b`).
				// Therefore we need to have different InputNames for them. Note that
				// we cannot just take the alias, as there would be the danger of
				// overriding an input stream with that same name, as in
				// `FROM x AS b, udsf('s') AS x`, and we should also use a string
				// that does not possibly conflict with any input name.
				// Note that `addTupleToBuffer` in defaultSelectExecutionPlan needs
				// to use that same method.
				InputName: fmt.Sprintf("%s/%s", rel.Name, alias),
			}); err != nil {
				return nil, err
			}
			bn.StopOnDisconnect(core.Inbound | core.Outbound)
			bn.RemoveOnStop()

		default:
			return nil, fmt.Errorf("input stream of type %s not implemented",
				rel.Type)
		}
	}
	dbox.(core.BoxNode).StopOnDisconnect(core.Inbound)
	removeNodes = false
	return dbox, nil
}

func (tb *TopologyBuilder) mkParamsMap(params []parser.SourceSinkParamAST) data.Map {
	paramsMap := make(data.Map, len(params))
	for _, kv := range params {
		paramsMap[string(kv.Key)] = kv.Value
	}
	return paramsMap
}

type chanSink struct {
	m      sync.RWMutex
	ch     chan *core.Tuple
	closed bool
}

func newChanSink() (*chanSink, <-chan *core.Tuple) {
	c := make(chan *core.Tuple)
	return &chanSink{
		ch: c,
	}, c
}

func (s *chanSink) Write(ctx *core.Context, t *core.Tuple) error {
	s.m.RLock()
	defer s.m.RUnlock()
	if s.closed {
		return errors.New("the sink has already been closed")
	}
	s.ch <- t
	return nil
}

func (s *chanSink) Close(ctx *core.Context) error {
	go func() {
		// Because Write might be blocked in s.ch <- t, this goroutine vacuums
		// tuples from the chan to unblock it and release the lock. Reading on
		// a closed chan is safe.
		for _ = range s.ch {
		}
	}()

	s.m.Lock()
	defer s.m.Unlock()
	if s.closed {
		return nil
	}
	s.closed = true
	close(s.ch)
	return nil
}

// AddSelectStmt creates nodes handling a SELECT statement in the topology.
// It returns the Sink node and the channel tied to it, the chan receiving
// tuples from the Sink, and an error if happens. The caller must stop the
// Sink node once it get unnecessary.
func (tb *TopologyBuilder) AddSelectStmt(stmt *parser.SelectStmt) (core.SinkNode, <-chan *core.Tuple, error) {
	sink, ch := newChanSink()
	temporaryName := fmt.Sprintf("sensorbee_tmp_select_sink_%v", topologyBuilderNextTemporaryID())
	sn, err := tb.topology.AddSink(temporaryName, sink, &core.SinkConfig{
		RemoveOnStop: true,
	})
	if err != nil {
		sink.Close(tb.topology.Context())
		return nil, nil, err
	}

	_, err = tb.AddStmt(parser.InsertIntoSelectStmt{
		Sink:       parser.StreamIdentifier(temporaryName),
		SelectStmt: *stmt,
	})
	if err != nil {
		if err := sn.Stop(); err != nil {
			tb.topology.Context().ErrLog(err).WithField("node_type", core.NTSink).
				WithField("node_name", temporaryName).Error("Cannot stop the temporary sink")
		}
		return nil, nil, err
	}
	sn.StopOnDisconnect()
	return sn, ch, nil
}
