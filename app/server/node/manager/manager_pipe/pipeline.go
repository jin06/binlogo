package manager_pipe

import (
	"context"
	"github.com/jin06/binlogo/pkg/store/dao/dao_sche"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/store/model/scheduler"
	"github.com/jin06/binlogo/pkg/watcher/scheduler_binding"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Manager struct {
	mapping        map[string]bool
	mappingIns     map[string]*instance
	bindingWatcher *scheduler_binding.BindingWatcher
	node           *node.Node
	mutex          sync.Mutex
}

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

func (m *Manager) Run(ctx context.Context) {
	go func() {
		if err := m.scanPipelines(nil); err != nil {
			logrus.Error(err)
		}
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			case <-time.Tick(time.Second * 5):
				{
					if err := m.scanPipelines(nil); err != nil {
						logrus.Error(err)
					}
				}
			case <-time.Tick(time.Second * 1):
				{
					m.dispatch()
				}
			}
		}
	}()
	return
}

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

	for pName, _ := range m.mapping {
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
			}
			go m.mappingIns[pName].start()
		}
		if !shouldRun {
			if isExist {
				m.mappingIns[pName].stop()
				delete(m.mappingIns, pName)
			}
		}
	}

}