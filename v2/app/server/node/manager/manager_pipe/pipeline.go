package manager_pipe

import (
	"context"
	"sync"
	"time"

	"github.com/jin06/binlogo/v2/pkg/store/dao/dao_sche"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
	"github.com/jin06/binlogo/v2/pkg/store/model/scheduler"
	"github.com/sirupsen/logrus"
)

// Manager is pipeline instance manager.
// If there is a pipeline bound to the current node, run an instance
// will also stop inbound instance on the current node
type Manager struct {
	mapping    map[string]bool
	mappingIns map[string]*instance
	node       *node.Node
	mutex      sync.Mutex
	stopped    chan struct{}
	closeOnce  sync.Once
}

// New returns a new Manager
func New(n *node.Node) (m *Manager) {
	m = &Manager{
		mapping:    map[string]bool{},
		mappingIns: map[string]*instance{},
		mutex:      sync.Mutex{},
		stopped:    make(chan struct{}),
	}
	m.node = n
	logrus.Debug("new pipeline manager ")
	logrus.Debug("nodeName:", n.Name)
	return
}

// Run start woking
func (m *Manager) Run(ctx context.Context) {
	var err error
	stx, cancel := context.WithCancel(ctx)
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln("pipeline manager panic, ", r)
		}
		cancel()
		m.close()
	}()
	if err = m.scanPipelines(nil); err != nil {
		logrus.Error(err)
	}
	ticker := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		case <-ticker.C:
			{
				if serr := m.scanPipelines(nil); serr != nil {
					logrus.Error(serr)
				}
				m.dispatch(stx)
			}
		}
	}
}

// scanPipelines scan pipeline bind, find pipelines that should run in this node
func (m *Manager) scanPipelines(pb *scheduler.PipelineBind) (err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if pb == nil {
		pb, err = dao_sche.GetPipelineBind()
		if err != nil {
			return
		}
	}
	for pipeName, nodeName := range pb.Bindings {
		if nodeName == m.node.Name {
			m.mapping[pipeName] = true
			continue
		}
		if nodeName != m.node.Name {
			m.mapping[pipeName] = false
		}
	}

	for pName := range m.mapping {
		if _, ok := pb.Bindings[pName]; !ok {
			m.mapping[pName] = false
		}
	}
	return
}

func (m *Manager) dispatch(ctx context.Context) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for pName, shouldRun := range m.mapping {
		//logrus.Debug(pName, shouldRun)
		ins, isExist := m.mappingIns[pName]
		if shouldRun {
			if !isExist || ins.exit {
				m.mappingIns[pName] = newInstance(pName, m.node.Name)
				go func() {
					m.mappingIns[pName].start(ctx)
				}()
				<-m.mappingIns[pName].started
			}
		}
		if !shouldRun {
			if isExist {
				m.mappingIns[pName].stop()
				delete(m.mappingIns, pName)
			}
		}
	}
}

func (m *Manager) close() {
	m.closeOnce.Do(func() {
		close(m.stopped)
	})
}

func (m *Manager) Stopped() chan struct{} {
	return m.stopped
}
