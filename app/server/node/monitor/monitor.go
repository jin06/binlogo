package monitor

import (
	"context"
	"github.com/jin06/binlogo/pkg/store/dao/dao_cluster"
	"github.com/jin06/binlogo/pkg/store/dao/dao_node"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
	"github.com/jin06/binlogo/pkg/watcher"
	"github.com/jin06/binlogo/pkg/watcher/node"
	"github.com/jin06/binlogo/pkg/watcher/pipeline"
	"sync"
)

type Monitor struct {
	status          string
	lock            sync.Mutex
	cancel          context.CancelFunc
	nodeWatcher     *watcher.General
	pipeWatcher     *watcher.General
	registerWatcher *watcher.General
}

const STATUS_RUN = "run"
const STATUS_STOP = "stop"

func NewMonitor() (m *Monitor, err error) {
	m = &Monitor{
		status: STATUS_STOP,
	}
	//m.pipeWatcher, err = pipeline.New(dao_pipe.PipelinePrefix())
	//m.nodeWatcher, err = node.New(dao_node.NodePrefix())
	//m.registerWatcher, err = node.New(dao_cluster.RegisterPrefix())
	return
}

func (m *Monitor) init() (err error) {
	m.pipeWatcher, err = pipeline.New(dao_pipe.PipelinePrefix())
	if err != nil {
		return
	}
	m.nodeWatcher, err = node.New(dao_node.NodePrefix())
	if err != nil {
		return
	}
	m.registerWatcher, err = node.New(dao_cluster.RegisterPrefix())
	if err != nil {
		return
	}
	return
}

func (m *Monitor) Run(ctx context.Context) (err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.status == STATUS_RUN {
		return
	}
	err = m.init()
	if err != nil {
		return
	}
	myCtx, cancel := context.WithCancel(ctx)
	defer func() {
		if err != nil {
			cancel()
		} else {
			m.status = STATUS_RUN
		}
	}()
	m.cancel = cancel
	err = m.monitorNode(myCtx)
	if err != nil {
		return
	}
	//err = m.monitorRegister(nctx)
	//if err != nil {
	//	return
	//}
	err = m.monitorPipe(myCtx)
	if err != nil {
		return
	}
	err = m.monitorStatus(myCtx)
	if err != nil {
		return
	}

	//err = m.monitorRegister(nctx)
	//if err != nil {
	//	return
	//}
	return nil
}

func (m *Monitor) Stop(ctx context.Context) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.status == STATUS_STOP {
		//blog.Debug("Monitor is stopped, do nothing")
		return
	}
	m.cancel()
	m.status = STATUS_STOP
	return
}
