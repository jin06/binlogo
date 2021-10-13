package scheduler

import (
	"context"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/jin06/binlogo/pkg/store/model"
	pipeline2 "github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/jin06/binlogo/pkg/store/model/scheduler"
	"github.com/jin06/binlogo/pkg/watcher/node"
	"github.com/jin06/binlogo/pkg/watcher/pipeline"
)

type Monitor struct {
	etcd              *etcd.ETCD
	pipelineWatcher   *pipeline.PipelineWatcher
	nodeWatcher       *node.NodeWatcher
	notBindPipelineCh chan *pipeline2.Pipeline
}

func newMonitor() (m *Monitor, err error) {
	m = &Monitor{
		etcd: etcd.E,
	}
	m.notBindPipelineCh = make(chan *pipeline2.Pipeline, 10000)
	prefix := etcd.Prefix()
	m.pipelineWatcher, err = pipeline.New(prefix + "/pipeline")
	m.nodeWatcher, err = node.New(prefix + "/nodes")
	return
}

func (m *Monitor) run(ctx context.Context) {
	blog.Debug("monitor run")
	go func() {
		m.pipelineWatcher.WatchList(ctx)
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			case p := <-m.pipelineWatcher.Queue:
				{
					blog.Debug("monitor queue \n", p)
					if ok, err := m.isBind(p); err != nil {
						if !ok {
							m.notBindPipelineCh <- p
						}
					}
				}
			}
		}
	}()
	return
}

func (m *Monitor) allPipelines() (res []*pipeline2.Pipeline, err error) {
	res = []*pipeline2.Pipeline{}
	resp, err := etcd.E.List("")
	for _, v := range resp {
		if val, ok := v.(*pipeline2.Pipeline); ok {
			res = append(res, val)
		}
	}
	return
}

func (m *Monitor) allNodes() (res []*model.Node, err error) {
	res = []*model.Node{}
	resp, err := etcd.E.List("")
	for _, v := range resp {
		if val, ok := v.(*model.Node); ok {
			res = append(res, val)
		}
	}
	return
}

func (m *Monitor) pipelineBind() (pb *scheduler.PipelineBind, err error) {
	pb = &scheduler.PipelineBind{}
	_, err = etcd.E.Get(pb)
	return
}

func (m *Monitor) isBind(p *pipeline2.Pipeline) (is bool, err error) {
	pb, err := m.pipelineBind()
	if err != nil {
		return
	}
	if len(pb.Bindings) == 0 {
		return
	}
	if _, ok := pb.Bindings[p.Name]; ok {
		is = true
		return
	}
	return
}
