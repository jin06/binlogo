package monitor

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

// Monitor monitor the operation of pipelines, nodes and other resources
type Monitor struct {
	status   string
	lock     sync.Mutex
	stopOnce sync.Once
	stopping chan struct{}
	Exit     bool
}

// NewMonitor returns a new Monitor
func NewMonitor() (m *Monitor) {
	m = &Monitor{
		stopping: make(chan struct{}),
	}
	//m.pipeWatcher, err = pipeline.New(dao.PipelinePrefix())
	//m.nodeWatcher, err = node.New(dao_node.NodePrefix())
	//m.registerWatcher, err = node.New(dao_cluster.RegisterPrefix())
	return
}

func (m *Monitor) init() (err error) {
	return
}

// Run start working
func (m *Monitor) Run(ctx context.Context) (err error) {
	logrus.Info("monitor run ")
	defer func() {
		m.Exit = true
	}()
	defer func() {
		if err != nil {
			logrus.WithError(err).Errorln("monitor stop")
		} else {
			logrus.Info("monitor stop")

		}
	}()
	if err = m.init(); err != nil {
		return
	}
	stx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		err = m.monitorNode(stx)
		m.stop()
	}()
	go func() {
		err = m.monitorPipe(stx)
		m.stop()
	}()
	go func() {
		err = m.monitorStatus(stx)
		m.stop()
	}()
	select {
	case <-ctx.Done():
	case <-m.stopping:
	}
	return nil
}

func (m *Monitor) Stop() {
	m.stop()
}

func (m *Monitor) stop() {
	m.stopOnce.Do(func() {
		close(m.stopping)
	})
}
