package manager_pipe

import (
	"context"
	"sync"
	"time"

	"github.com/jin06/binlogo/pkg/store/dao/dao_sche"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/store/model/scheduler"
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
	ctx        context.Context
}

// New returns a new Manager
func New(n *node.Node) (m *Manager) {
	m = &Manager{
		mapping:    map[string]bool{},
		mappingIns: map[string]*instance{},
		mutex:      sync.Mutex{},
	}
	m.node = n
	logrus.Debug("new pipeline manager ")
	logrus.Debug("nodeName:", n.Name)
	return
}

// Run start woking
func (m *Manager) Run(ctx context.Context) {
	myCtx, cancel := context.WithCancel(ctx)
	m.ctx = myCtx
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorln("election manager panic, ", r)
			}
			cancel()
		}()
		if err := m.scanPipelines(nil); err != nil {
			logrus.Error(err)
		}
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			case <-time.Tick(time.Second * 1):
				{
					if err := m.scanPipelines(nil); err != nil {
						logrus.Error(err)
					}
					m.dispatch()
				}
				//case <-time.Tick(time.Second * 1):
				//	{
				//		m.dispatch()
				//	}
			}
		}
	}()
	return
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

func (m *Manager) dispatch() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for pName, shouldRun := range m.mapping {
		_, isExist := m.mappingIns[pName]
		if shouldRun {
			if !isExist {
				newIns := newInstance(pName, m.node.Name)
				m.mappingIns[pName] = newIns
				//go m.mappingIns[pName].start(m.ctx)
				go func() {
					ins := m.mappingIns[pName]
					runErr := ins.start(m.ctx)
					logrus.WithField("pipeline name", pName).WithField("error", runErr).Warn("pipeline exist")
					m.mutex.Lock()
					delete(m.mappingIns, pName)
					m.mutex.Unlock()
				}()
			}
		}
		if !shouldRun {
			if isExist {
				m.mappingIns[pName].stop()
				//delete(m.mappingIns, pName)
			}
		}
	}
}

// Context returns Manager ctx
func (m *Manager) Context() context.Context {
	return m.ctx
}
