package core

import (
	"fmt"
	"sync"
)

type defaultDynamicTopology struct {
	ctx  *Context
	name string

	// nodeMutex protects sources, boxes, and sinks from being modified
	// concurrently. DO NOT acquire this lock after locking stateMutex
	// (opposite is fine).
	nodeMutex sync.RWMutex
	sources   map[string]*defaultDynamicSourceNode
	boxes     map[string]*defaultDynamicBoxNode
	sinks     map[string]*defaultDynamicSinkNode

	state      *topologyStateHolder
	stateMutex sync.Mutex

	// TODO: support lazy invocation of GenerateStream (call it when the first
	// destination is added or a Sink is indirectly connected). Maybe graph
	// management is required.
}

// NewDefaultDynamicTopology creates a dynamic topology having a simple graph
// structure.
func NewDefaultDynamicTopology(ctx *Context, name string) DynamicTopology {
	// TODO: validate name

	t := &defaultDynamicTopology{
		ctx:  ctx,
		name: name,

		sources: map[string]*defaultDynamicSourceNode{},
		boxes:   map[string]*defaultDynamicBoxNode{},
		sinks:   map[string]*defaultDynamicSinkNode{},
	}
	t.state = newTopologyStateHolder(&t.stateMutex)
	t.state.state = TSRunning // A dynamic topology is running by default.
	return t
}

// AddSource adds a source to the topology. It will asynchronously call source's
// GenerateStream and returns after the source started running. GenerateStream
// could be called lazily to avoid unnecessary computation. GenerateStream and
// Stop might also be called when this method returns an error. The caller must
// not call GenerateStream or Stop of the source.
func (t *defaultDynamicTopology) AddSource(name string, s Source, config *DynamicSourceConfig) (DynamicSourceNode, error) {
	// TODO: validate the name

	// This method assumes adding a Source having a duplicated name is rare.
	// Under this assumption, acquiring wlock without checking the existence
	// of the name with rlock doesn't degrade the performance.
	t.nodeMutex.Lock()
	defer t.nodeMutex.Unlock()

	// t.state is set to TSStopped while t.nodeMutex is locked. Therefore,
	// a source can safely be added when the state is TSRunning or TSPaused.
	//
	// The lock above doesn't protect t.state from being set to TSStopping.
	// So, some goroutine can go through this if block and add a source to
	// the topology when the state is TSStopping. However, adding a source
	// while the state is TSStopping is safe although the source will get
	// removed just after it's added to the topology.
	if t.state.Get() >= TSStopping {
		return nil, fmt.Errorf("the topology is already stopped")
	}

	ds := &defaultDynamicSourceNode{
		defaultDynamicNode: newDefaultDynamicNode(t, name),
		source:             s,
		dsts:               newDynamicDataDestinations(name),
	}
	if err := t.checkNodeNameDuplication(name); err != nil {
		if err := ds.Stop(); err != nil { // The same variable name is intentionally used.
			t.ctx.Logger.Log(Error, "Cannot stop the source '%v': %v", name, err)
		}
		return nil, err
	}
	t.sources[name] = ds

	go func() {
		// TODO: Support lazy invocation
		if err := ds.run(); err != nil {
			t.ctx.Logger.Log(Error, "Cannot generate a stream from the source '%v': %v", name, err)
		}
	}()
	ds.state.Wait(TSRunning)
	return ds, nil
}

// TODO: Add method to validate a node name

// checkNodeNameDuplication checks if the given name is unique in the topology.
// This method doesn't acquire the lock and it's the caller's responsibility
// to do it before calling this method.
func (t *defaultDynamicTopology) checkNodeNameDuplication(name string) error {
	if _, ok := t.sources[name]; ok {
		return fmt.Errorf("the name is already used by a source: %v", name)
	}
	if _, ok := t.boxes[name]; ok {
		return fmt.Errorf("the name is already used by a box: %v", name)
	}
	if _, ok := t.sinks[name]; ok {
		return fmt.Errorf("the name is already used by a sink: %v", name)
	}
	return nil
}

func (t *defaultDynamicTopology) AddBox(name string, b Box, config *DynamicBoxConfig) (DynamicBoxNode, error) {
	t.nodeMutex.Lock()
	defer t.nodeMutex.Unlock()
	if t.state.Get() >= TSStopping {
		return nil, fmt.Errorf("the topology is already stopped")
	}

	if err := t.checkNodeNameDuplication(name); err != nil {
		return nil, err
	}

	if sb, ok := b.(StatefulBox); ok {
		if err := sb.Init(t.ctx); err != nil {
			return nil, err
		}
	}

	db := &defaultDynamicBoxNode{
		defaultDynamicNode: newDefaultDynamicNode(t, name),
		srcs:               newDynamicDataSources(name),
		box:                b,
		dsts:               newDynamicDataDestinations(name),
	}
	t.boxes[name] = db

	go func() {
		if err := db.run(); err != nil {
			t.ctx.Logger.Log(Error, "Box '%v' in topology '%v' failed: %v", db.name, t.name, err)
		}
	}()
	db.state.Wait(TSRunning)
	return db, nil
}

func (t *defaultDynamicTopology) AddSink(name string, s Sink, config *DynamicSinkConfig) (DynamicSinkNode, error) {
	t.nodeMutex.Lock()
	defer t.nodeMutex.Unlock()
	if t.state.Get() >= TSStopping {
		return nil, fmt.Errorf("the topology is already stopped")
	}

	if err := t.checkNodeNameDuplication(name); err != nil {
		return nil, err
	}

	ds := &defaultDynamicSinkNode{
		defaultDynamicNode: newDefaultDynamicNode(t, name),
		srcs:               newDynamicDataSources(name),
		sink:               s,
	}
	t.sinks[name] = ds

	go func() {
		if err := ds.run(); err != nil {
			t.ctx.Logger.Log(Error, "Sink '%v' in topology '%v' failed: %v", ds.name, t.name, err)
		}
	}()
	ds.state.Wait(TSRunning)
	return ds, nil
}

func (t *defaultDynamicTopology) Stop() error {
	stopped, err := func() (bool, error) {
		t.stateMutex.Lock()
		defer t.stateMutex.Unlock()
		switch t.state.state {
		case TSRunning, TSPaused:
			t.state.setWithoutLock(TSStopping)
			return false, nil

		case TSStopping:
			t.state.waitWithoutLock(TSStopped)
			return true, nil

		case TSStopped:
			return true, nil

		default: // including TSInitialized and TSStarting
			return false, fmt.Errorf("the dynamic topology has an invalid state: %v", t.state.state)
		}
	}()
	if err != nil {
		return err
	} else if stopped {
		return nil
	}

	t.nodeMutex.Lock()
	defer t.nodeMutex.Unlock()

	var lastErr error
	for name, src := range t.sources {
		// TODO: this could be run concurrently
		func() {
			defer func() {
				if e := recover(); e != nil {
					t.ctx.Logger.Log(Error, "Cannot stop source '%v' due to panic: %v", name, e)
					src.dsts.Close(t.ctx)

					if err, ok := e.(error); ok {
						lastErr = err
					} else {
						lastErr = fmt.Errorf("source '%v' panicked while being stopped: %v", name, e)
					}
				}
			}()

			if err := src.Stop(); err != nil {
				src.dsts.Close(t.ctx)
				lastErr = err
				t.ctx.Logger.Log(Error, "Cannot stop source '%v': %v", name, err)
			}
		}()
	}

	var wg sync.WaitGroup
	for _, b := range t.boxes {
		b.srcs.enableGracefulStop()
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.state.Wait(TSStopped)
		}()
	}

	for _, s := range t.sinks {
		s.srcs.enableGracefulStop()
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.state.Wait(TSStopped)
		}()
	}
	wg.Wait()

	t.sources = nil
	t.boxes = nil
	t.sinks = nil
	t.state.Set(TSStopped)
	return nil
}

// TODO: Add Remove method to remove a node from the topology
// TODO: Add method to clean up (possibly indirectly) stopped nodes
// TODO: Add Node(name), Source(name), Box(name), Sink(name) methods to retrieve specific nodes in the topology

func (t *defaultDynamicTopology) Node(name string) (DynamicNode, error) {
	t.nodeMutex.RLock()
	defer t.nodeMutex.RUnlock()
	return t.nodeWithoutLock(name)
}

func (t *defaultDynamicTopology) nodeWithoutLock(name string) (DynamicNode, error) {
	if s, ok := t.sources[name]; ok {
		return s, nil
	}
	if b, ok := t.boxes[name]; ok {
		return b, nil
	}
	if s, ok := t.sinks[name]; ok {
		return s, nil
	}
	return nil, fmt.Errorf("node '%v' was not found", name)
}

type dynamicDataSource interface {
	Name() string
	destinations() *dynamicDataDestinations
}

func (t *defaultDynamicTopology) dataSource(nodeName string) (dynamicDataSource, error) {
	t.nodeMutex.RLock()
	defer t.nodeMutex.RUnlock()

	if s, ok := t.sources[nodeName]; ok {
		return s, nil
	}
	if b, ok := t.boxes[nodeName]; ok {
		return b, nil
	}
	return nil, fmt.Errorf("no such data source node: %v", nodeName)
}

type defaultDynamicNode struct {
	topology   *defaultDynamicTopology
	name       string
	state      *topologyStateHolder
	stateMutex sync.Mutex
}

func newDefaultDynamicNode(t *defaultDynamicTopology, name string) *defaultDynamicNode {
	dn := &defaultDynamicNode{
		topology: t,
		name:     name,
	}
	dn.state = newTopologyStateHolder(&dn.stateMutex)
	return dn
}

func (dn *defaultDynamicNode) Name() string {
	return dn.name
}

func (dn *defaultDynamicNode) State() TopologyStateHolder {
	return dn.state
}

func (dn *defaultDynamicNode) checkAndPrepareRunState(nodeType string) error {
	dn.stateMutex.Lock()
	defer dn.stateMutex.Unlock()
	switch s := dn.state.getWithoutLock(); s {
	case TSInitialized:
		dn.state.setWithoutLock(TSStarting)
		return nil

	case TSStarting:
		dn.state.waitWithoutLock(TSRunning)
		fallthrough

	case TSRunning, TSPaused:
		return fmt.Errorf("%v '%v' is already running", nodeType, dn.name)

	case TSStopping:
		dn.state.waitWithoutLock(TSStopped)
		fallthrough

	case TSStopped:
		return fmt.Errorf("%v '%v' is already stopped", nodeType, dn.name)

	default:
		return fmt.Errorf("%v '%v' has an invalid state: %v", nodeType, dn.name, s)
	}
}

// checkAndPrepareStopState check the current state of the node and returns if
// the node can be stopped or is already stopped.
func (dn *defaultDynamicNode) checkAndPrepareStopState(nodeType string) (stopped bool, err error) {
	dn.stateMutex.Lock()
	defer dn.stateMutex.Unlock()
	for {
		switch s := dn.state.getWithoutLock(); s {
		case TSInitialized:
			dn.state.setWithoutLock(TSStopped)
			return true, nil

		case TSStarting:
			dn.state.waitWithoutLock(TSRunning)
			// try again in the next iteration

		case TSRunning, TSPaused:
			return false, nil

		case TSStopping:
			dn.state.waitWithoutLock(TSStopped)
			fallthrough

		case TSStopped:
			return true, nil

		default:
			return false, fmt.Errorf("%v '%v' has an invalid state: %v", nodeType, dn.name, s)
		}
	}
}
