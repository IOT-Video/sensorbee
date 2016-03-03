package bql

import (
	"errors"
	"fmt"
	"gopkg.in/sensorbee/sensorbee.v0/bql/execution"
	"gopkg.in/sensorbee/sensorbee.v0/bql/parser"
	"gopkg.in/sensorbee/sensorbee.v0/bql/udf"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"math"
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
	UDSStorage     udf.UDSStorage
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
	if err := srcs.Register("edge_statuses", createEdgeStatusSourceCreator(t)); err != nil {
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
		UDSStorage:     udf.NewInMemoryUDSStorage(),
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

		// check if we know this type of source
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

	case parser.CreateStreamAsSelectUnionStmt:
		// idea: create an intermediate box for each SELECT substatement,
		// then connect them with a simple forwarder box
		names := make([]string, 0, len(stmt.Selects))
		nodes := make([]core.BoxNode, 0, len(stmt.Selects))
		removeTmpNodes := func() {
			for _, name := range names {
				tb.topology.Remove(name)
			}
		}
		for _, selStmt := range stmt.Selects {
			// create a stream with a generated name and recurse
			tmpName := fmt.Sprintf("sensorbee_tmp_%v", topologyBuilderNextTemporaryID())
			tmpStmt := parser.CreateStreamAsSelectStmt{
				parser.StreamIdentifier(tmpName),
				selStmt,
			}
			box, err := tb.AddStmt(tmpStmt)
			if err != nil {
				removeTmpNodes()
				return nil, err
			}
			names = append(names, tmpName)
			nodes = append(nodes, box.(core.BoxNode))
		}
		// simple forwarder box
		forwardBox := core.BoxFunc(func(ctx *core.Context, t *core.Tuple, w core.Writer) error {
			return w.Write(ctx, t)
		})
		node, err := tb.topology.AddBox(string(stmt.Name), forwardBox, nil)
		if err != nil {
			removeTmpNodes()
			return nil, err
		}
		// connect inputs
		for _, name := range names {
			if err := node.Input(name, nil); err != nil {
				removeTmpNodes()
				return nil, err
			}
		}
		for _, node := range nodes {
			node.StopOnDisconnect(core.Inbound | core.Outbound)
			node.RemoveOnStop()
		}
		node.StopOnDisconnect(core.Inbound)
		node.RemoveOnStop()
		return node, nil

	case parser.CreateSinkStmt:
		// load params into map for faster access
		paramsMap := tb.mkParamsMap(stmt.Params)

		// check if we know this type of sink
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
		if err := ctx.SharedStates.Add(string(stmt.Name), string(stmt.Type), s); err != nil {
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
		return nil, u.Update(ctx, tb.mkParamsMap(stmt.Params))

	case parser.SaveStateStmt:
		return nil, tb.saveState(string(stmt.Name), stmt.Tag)

	case parser.LoadStateStmt:
		_, err := tb.loadState(string(stmt.Type), string(stmt.Name), stmt.Tag, tb.mkParamsMap(stmt.Params))
		return nil, err

	case parser.LoadStateOrCreateStmt:
		shouldCreate, err := tb.loadState(string(stmt.Type), string(stmt.Name), stmt.Tag, tb.mkParamsMap(stmt.LoadSpecs.Params))
		if shouldCreate {
			c := parser.CreateStateStmt{}
			c.Type = stmt.Type
			c.Name = stmt.Name
			c.Params = stmt.CreateSpecs.Params
			return tb.AddStmt(c)
		}
		return nil, err

	case parser.UpdateSourceStmt:
		src, err := tb.topology.Source(string(stmt.Name))
		if err != nil {
			return nil, err
		}

		u, ok := src.Source().(core.Updater)
		if !ok {
			return nil, fmt.Errorf("%s cannot be updated", string(stmt.Name))
		}
		return nil, u.Update(tb.topology.Context(), tb.mkParamsMap(stmt.Params))

	case parser.UpdateSinkStmt:
		sink, err := tb.topology.Sink(string(stmt.Name))
		if err != nil {
			return nil, err
		}

		u, ok := sink.Sink().(core.Updater)
		if !ok {
			return nil, fmt.Errorf("%s cannot be updated", string(stmt.Name))
		}
		return nil, u.Update(tb.topology.Context(), tb.mkParamsMap(stmt.Params))

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

// udsfBox is a core.Box which runs a UDSF in the stream mode.
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

// udsfSource is a core.Source which runs a UDSF in the source mode.
type udsfSource struct {
	f       udf.UDSF
	stopped core.AtomicFlag
}

var (
	_ core.Source = &udsfSource{}
)

func newUDSFSource(f udf.UDSF) *udsfSource {
	return &udsfSource{
		f: f,
	}
}

func (s *udsfSource) GenerateStream(ctx *core.Context, w core.Writer) error {
	// In the source mode, UDSF.Process is only called once. It can generate
	// as many tuples as it wants.
	return s.f.Process(ctx, core.NewTuple(data.Map{"b": data.True}),
		core.WriterFunc(func(ctx *core.Context, t *core.Tuple) error {
			if s.stopped.Enabled() {
				return core.ErrSourceStopped
			}
			return w.Write(ctx, t)
		}))
}

func (s *udsfSource) Stop(ctx *core.Context) error {
	s.stopped.Set(true)
	return s.f.Terminate(ctx)
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
	// provide a function to the BQL box to remove itself from the topology
	box.removeMe = func() { go tb.topology.Remove(outName) }

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
	var pausedSources []core.SourceNode
	for _, rel := range stmt.Select.Relations {
		switch rel.Type {
		case parser.ActualStream:
			if connected[rel.Name] {
				// this is a self-join (FROM x [RANGE ...] AS a, x [RANGE ...] AS b)
				// and we already have connected x to this box before
				continue
			}
			conf := &core.BoxInputConfig{
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
			}
			// set capacity of input pipe
			if rel.Capacity != parser.UnspecifiedCapacity {
				if rel.Capacity > math.MaxInt32 {
					return nil, fmt.Errorf("specified buffer capacity %d is too large",
						rel.Capacity)
				} else if rel.Capacity < 0 {
					// the parser should not allow this to happen, actually
					return nil, fmt.Errorf("specified buffer capacity %d must not be negative",
						rel.Capacity)
				}
				conf.Capacity = int(rel.Capacity)
			}
			// set drop mode for box
			if rel.Shedding == parser.DropOldest {
				conf.DropMode = core.DropOldest
			} else if rel.Shedding == parser.DropNewest {
				conf.DropMode = core.DropLatest
			} else if rel.Shedding == parser.Wait {
				conf.DropMode = core.DropNone
			}
			if err := dbox.Input(rel.Name, conf); err != nil {
				return nil, err
			}
			connected[rel.Name] = true

		case parser.UDSFStream:
			sn, name, err := tb.setUpUDSFStream(dbox, &rel)
			if err != nil {
				return nil, err
			}
			temporaryNodes = append(temporaryNodes, name)
			if sn != nil {
				pausedSources = append(pausedSources, sn)
			}

		default:
			return nil, fmt.Errorf("input stream of type %s not implemented",
				rel.Type)
		}
	}
	dbox.(core.BoxNode).StopOnDisconnect(core.Inbound)

	// Resume all UDSFs running in the source mode as fairly as possible.
	for _, sn := range pausedSources {
		if err := sn.Resume(); err != nil {
			return nil, err
		}
	}
	removeNodes = false
	return dbox, nil
}

// setUpUDSFStream creates a Source or a Box from a UDSF. When it creates a
// Source, it will return the corresponding core.SourceNode of it. Otherwise,
// it returns nil for core.SourceNode. It also returns the temporary name of
// the UDSF node.
func (tb *TopologyBuilder) setUpUDSFStream(subsequentBox core.BoxNode, rel *parser.AliasedStreamWindowAST) (core.SourceNode, string, error) {
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
			return nil, "", err
		}
		params[i] = p
	}

	udsfc, err := tb.UDSFCreators.Lookup(rel.Name, len(params))
	if err != nil {
		return nil, "", err
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
		return nil, "", err
	}

	temporaryName := fmt.Sprintf("sensorbee_tmp_udsf_%v", topologyBuilderNextTemporaryID())
	addInput := func() error {
		alias := rel.Alias
		if alias == "" {
			alias = rel.Name
		}
		conf := &core.BoxInputConfig{
			// As opposed to actual streams, for `udsf("s") AS a, udsf("s") AS b`,
			// there will be *multiple* boxes and we will have one connection to
			// udsf 1 (the one aliased to `a`) and udsf 2 (the one aliased to `b`).
			// Therefore we need to have different InputNames for them. Note that
			// we cannot just take the alias, as there would be the danger of
			// overriding an input stream with that same name, as in
			// `FROM x AS b, udsf("s") AS x`, and we should also use a string
			// that does not possibly conflict with any input name.
			// Note that `addTupleToBuffer` in defaultSelectExecutionPlan needs
			// to use that same method.
			InputName: fmt.Sprintf("%s/%s", rel.Name, alias),
		}
		// set capacity of input pipe
		if rel.Capacity != parser.UnspecifiedCapacity {
			if rel.Capacity > math.MaxInt32 {
				return fmt.Errorf("specified buffer capacity %d is too large", rel.Capacity)
			} else if rel.Capacity < 0 {
				// the parser should not allow this to happen, actually
				return fmt.Errorf("specified buffer capacity %d must not be negative", rel.Capacity)
			}
			conf.Capacity = int(rel.Capacity)
		}
		// set drop mode for box
		if rel.Shedding == parser.DropOldest {
			conf.DropMode = core.DropOldest
		} else if rel.Shedding == parser.DropNewest {
			conf.DropMode = core.DropLatest
		} else if rel.Shedding == parser.Wait {
			conf.DropMode = core.DropNone
		}
		return subsequentBox.Input(temporaryName, conf)
	}

	if len(decl.ListInputs()) == 0 { // Source mode
		sn, err := tb.topology.AddSource(temporaryName, newUDSFSource(udsf), &core.SourceConfig{
			PausedOnStartup: true,
		})
		if err != nil {
			return nil, "", err
		}
		if err := addInput(); err != nil {
			return nil, "", err
		}
		sn.StopOnDisconnect()
		sn.RemoveOnStop()
		// The source will be resumed by the caller.
		return sn, temporaryName, nil
	}

	bn, err := tb.topology.AddBox(temporaryName, newUDSFBox(udsf), &core.BoxConfig{
	// TODO: add information of the statement
	})
	if err != nil {
		return nil, "", err
	}
	for input, config := range decl.ListInputs() {
		if err := bn.Input(input, &core.BoxInputConfig{
			InputName: config.InputName,
		}); err != nil {
			return nil, "", err
		}
	}
	if err := addInput(); err != nil {
		return nil, "", err
	}
	bn.StopOnDisconnect(core.Inbound | core.Outbound)
	bn.RemoveOnStop()
	return nil, temporaryName, nil
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
	// wrap this in a UNION statement
	tmpStmt := parser.SelectUnionStmt{[]parser.SelectStmt{*stmt}}
	return tb.AddSelectUnionStmt(&tmpStmt)
}

// AddSelectUnionStmt creates nodes handling a SELECT ... UNION ALL statement
// in the topology. It returns the Sink node and the channel tied to it, the
// chan receiving tuples from the Sink, and an error if happens. The caller must
// stop the Sink node once it get unnecessary.
func (tb *TopologyBuilder) AddSelectUnionStmt(stmts *parser.SelectUnionStmt) (core.SinkNode, <-chan *core.Tuple, error) {
	sink, ch := newChanSink()
	tmpUnionNodeName := fmt.Sprintf("sensorbee_tmp_select_sink_%v", topologyBuilderNextTemporaryID())
	sn, err := tb.topology.AddSink(tmpUnionNodeName, sink, nil)
	if err != nil {
		sink.Close(tb.topology.Context())
		return nil, nil, err
	}

	names := make([]string, 0, len(stmts.Selects))
	for _, stmt := range stmts.Selects {
		// In an earlier version of the code, we used to process
		// this as an InsertIntoSelectStmt and insert into the
		// temporary node created above. InsertIntoSelectStmt
		// has been removed, so the corresponding code has
		// been moved here. Since error handling requires
		// quite a bit of cleanup, we wrap this functionality
		// into an inline function.

		createTmpBox := func(outputName string, stmt parser.SelectStmt) (core.Node, error) {
			// get the sink to add an input to
			sink, err := tb.topology.Sink(outputName)
			if err != nil {
				return nil, err
			}
			// construct an intermediate box doing the SELECT computation:
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
		}

		node, err := createTmpBox(tmpUnionNodeName, stmt)
		if err != nil {
			// clean up the already created nodes
			for _, name := range names {
				tb.topology.Remove(name)
			}
			if err := sn.Stop(); err != nil {
				tb.topology.Context().ErrLog(err).WithField("node_type", core.NTSink).
					WithField("node_name", tmpUnionNodeName).Error("Cannot stop the temporary sink")
			}
			tb.topology.Remove(tmpUnionNodeName)
			return nil, nil, err
		}
		names = append(names, node.Name())
	}
	sn.RemoveOnStop()
	sn.StopOnDisconnect()
	return sn, ch, nil
}

// RunEvalStmt evaluates the expression contained in the given EvalStmt
// and returns the evaluation result.
func (tb *TopologyBuilder) RunEvalStmt(stmt *parser.EvalStmt) (data.Value, error) {
	if stmt.Input == nil {
		// there is no ON clause, therefore our expression must
		// be foldable
		return execution.EvaluateFoldable(stmt.Expr, tb.Reg)
	}
	// if we arrive here, there was an ON clause given. first of all, we
	// must evaluate that ON expression
	inputData, err := execution.EvaluateFoldable(*stmt.Input, tb.Reg)
	if err != nil {
		return nil, err
	}
	// check that the expression we got is sane in this context
	usedRelations := stmt.Expr.ReferencedRelations()
	if len(usedRelations) > 1 || (len(usedRelations) == 1 && !usedRelations[""]) {
		return nil, fmt.Errorf("stream prefixes cannot be used inside EVAL")
	}
	expr := stmt.Expr.RenameReferencedRelation("", "input")
	// nest the data so that access via JSON path works properly
	inputRow := data.Map{"input": inputData}
	return execution.EvaluateOnInput(expr, inputRow, tb.Reg)
}

func (tb *TopologyBuilder) saveState(name, tag string) error {
	st, err := tb.topology.Context().SharedStates.Get(name)
	if err != nil {
		return err
	}
	s, ok := st.(core.SavableSharedState)
	if !ok {
		return fmt.Errorf("the state '%v-%v' cannot be saved", name, tag)
	}

	// Appropriate header information should be written by the storage.
	w, err := tb.UDSStorage.Save(tb.topology.Name(), name, tag)
	if err != nil {
		return err
	}
	shouldAbort := true
	defer func() {
		if shouldAbort {
			if err := w.Abort(); err != nil {
				tb.topology.Context().ErrLog(err).WithField("state_name", name).
					WithField("state_tag", tag).
					Error("saving the state panicked")
			}
		}
	}()

	if err := s.Save(tb.topology.Context(), w, data.Map{}); err != nil {
		return err
	}
	shouldAbort = false
	return w.Commit()
}

// loadState loads a state from the storage. It returns true when the state was
// not saved and LOAD STATE OR CREATE IF NOT SAVED should fall back to CREATE STATE.
func (tb *TopologyBuilder) loadState(typeName, name, tag string, params data.Map) (bool, error) {
	r, err := tb.UDSStorage.Load(tb.topology.Name(), name, tag)
	if err != nil {
		return core.IsNotExist(err), err
	}
	defer r.Close()

	c, err := tb.UDSCreators.Lookup(typeName)
	if err != nil {
		return false, err
	}
	loader, ok := c.(udf.UDSLoader)
	if !ok {
		return false, fmt.Errorf("the state '%v-%v' cannot be loaded", name, tag)
	}

	// If the state is loaded and it provides Load method, Load method will be
	// used instead of loader.LoadState.
	reg := tb.topology.Context().SharedStates
	s, err := reg.Get(name)
	if err != nil {
		// TODO: check if the error is "not found". Return only on other errors.
	} else if t, err := reg.Type(name); err != nil {
		return false, err
	} else if t != typeName {
		return false, fmt.Errorf("type name doesn't much to the current state's type")
	}

	if l, ok := s.(core.LoadableSharedState); ok {
		return false, l.Load(tb.topology.Context(), r, params)
	}

	newState, err := loader.LoadState(tb.topology.Context(), r, params)
	if err != nil {
		return false, err
	}
	prev, err := reg.Replace(name, typeName, newState)
	if err != nil {
		return false, err
	}
	if prev != nil {
		if err := prev.Terminate(tb.topology.Context()); err != nil {
			tb.topology.Context().ErrLog(err).WithField("state_name", name).
				Error("Cannot terminate the previous instance of the loaded state")
		}
	}
	return false, nil
}
