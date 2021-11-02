package monitor

import (
	"context"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/store/dao/dao_cluster"
	"github.com/jin06/binlogo/pkg/store/dao/dao_node"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
	"github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/jin06/binlogo/pkg/watcher/node"
	"github.com/jin06/binlogo/pkg/watcher/pipeline"
	"sync"
)

type Monitor struct {
	status          string
	lock            sync.Mutex
	cancel          context.CancelFunc
	nodeWatcher     *node.NodeWatcher
	pipeWatcher     *pipeline.PipelineWatcher
	registerWatcher *node.NodeWatcher
	etcd            *etcd.ETCD
}

const STATUS_RUN = "run"
const STATUS_STOP = "stop"

func NewMonitor() (m *Monitor, err error) {
	m = &Monitor{
		etcd: etcd.E,
		status: STATUS_STOP,
	}
	m.pipeWatcher, err = pipeline.New(dao_pipe.PipelinePrefix())
	m.nodeWatcher, err = node.New(dao_node.NodePrefix())
	m.registerWatcher, err = node.New(dao_cluster.RegisterPrefix())
	return
}

func (m *Monitor) Run(ctx context.Context) (err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.status == STATUS_RUN {
		blog.Debug("Monitor is running, do nothing")
		return
	}
	nctx, cancel := context.WithCancel(ctx)
	m.cancel = cancel
	m.monitorNode(nctx)
	m.monitorRegister(nctx)
	m.monitorPipe(nctx)
	m.status = STATUS_RUN

	return nil
}

func (m *Monitor) Stop(ctx context.Context) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.status == STATUS_STOP {
		blog.Debug("Monitor is stopped, do nothing")
		return
	}
	m.cancel()
	m.status = STATUS_STOP
	return
}
