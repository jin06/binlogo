package monitor

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

// Monitor monitor the operation of pipelines, nodes and other resources
type Monitor struct {
	closing      chan struct{}
	closeOnce    sync.Once
	closed       chan struct{}
	completeOnce sync.Once
	isClosed     bool
}

// NewMonitor returns a new Monitor
func NewMonitor() (m *Monitor) {
	m = &Monitor{
		closing: make(chan struct{}),
		closed:  make(chan struct{}),
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
	defer m.CompleteClose()
	defer m.Close()
	logrus.Info("monitor run ")
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
	go func() {
		err = m.monitorNode(ctx)
		m.Close()
	}()
	go func() {
		err = m.monitorPipe(ctx)
		m.Close()
	}()
	go func() {
		err = m.monitorStatus(ctx)
		m.Close()
	}()
	select {
	case <-ctx.Done():
	case <-m.closing:
	}
	return nil
}

func (m *Monitor) Close() {
	m.closeOnce.Do(func() {
		close(m.closing)
	})
}

func (m *Monitor) CompleteClose() {
	m.completeOnce.Do(func() {
		m.isClosed = true
		close(m.closed)
	})
}

func (m *Monitor) IsClosed() bool {
	return m.isClosed
}

func (m *Monitor) Closed() chan struct{} {
	return m.closed
}
