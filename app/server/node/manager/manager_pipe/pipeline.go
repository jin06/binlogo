package manager_pipe

import (
	"context"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/store/dao/dao_sche"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/store/model/scheduler"
	"github.com/jin06/binlogo/pkg/watcher/scheduler_binding"
	"time"
)

type Manager struct {
	mapping        map[string]*instance
	bindingWatcher *scheduler_binding.BindingWatcher
	node           *node.Node
}

func New(n *node.Node) (m *Manager) {
	m = &Manager{
		mapping: map[string]*instance{},
	}
	m.node = n
	blog.Debug("new pipeline manager ")
	blog.Debug("nodeName:", n.Name)
	return
}

func (m *Manager) Run(ctx context.Context) {
	go func() {
		if err := m.handlePipelineBind(nil); err != nil {
			blog.Error(err)
		}
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			case <-time.Tick(time.Second * 60):
				{
					if err := m.handlePipelineBind(nil); err != nil {
						blog.Error(err)
					}
				}
			}
		}
	}()
	return
}

func (m *Manager) handlePipelineBind(pb *scheduler.PipelineBind) (err error) {
	if pb == nil {
		pb, err = dao_sche.GetPipelineBind()
		if err != nil {
			return
		}
	}
	for k, v := range pb.Bindings {
		er := m.handle(k, v)
		if er != nil {
			blog.Error(er)
		}
	}
	return
}

func (m *Manager) handle(pipeName string, nodeName string) (err error) {
	blog.Debugln("manager_pipeline, ", pipeName, " nodeName", nodeName)
	if _, ok := m.mapping[pipeName]; ok {
		if nodeName != m.node.Name {
			err = m.removePipeline(pipeName)
		}
	} else {
		if nodeName == m.node.Name {
			blog.Debugln("manager_pipeline, addPipeline")
			err = m.addPipeline(pipeName)
		}
	}
	return
}

func (m *Manager) addPipeline(name string) (err error) {
	if _, ok := m.mapping[name]; ok {
		err = m.removePipeline(name)
		if err != nil {
			return
		}
	}
	var ins *instance
	ins, err = newInstance(name)
	if err != nil {
		return
	}
	err = ins.start()
	if err != nil {
		return
	}
	m.mapping[name] = ins
	return
}

func (m *Manager) removePipeline(name string) (err error) {
	if ins, ok := m.mapping[name]; ok {
		ins.p.Stop()
		err = ins.m.Unlock()
	}
	if err != nil {
		delete(m.mapping, name)
	}
	return
}
