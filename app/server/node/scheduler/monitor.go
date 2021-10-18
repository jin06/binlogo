package scheduler

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/store/dao"
	"github.com/jin06/binlogo/pkg/store/etcd"
	node2 "github.com/jin06/binlogo/pkg/store/model/node"
	pipeline2 "github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/jin06/binlogo/pkg/store/model/scheduler"
	"github.com/jin06/binlogo/pkg/watcher/node"
	"github.com/jin06/binlogo/pkg/watcher/pipeline"
	"time"
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
		m.nodeWatcher.WatchList(ctx)
		m.queueAllPipelines()
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			case p := <-m.pipelineWatcher.Queue:
				{
					blog.Debug("monitor queue \n", p)
					if ok, err := m.isBind(p); err == nil {
						if !ok {
							m.notBindPipelineCh <- p
						}
					} else {
						blog.Debug("check bind error", err)
					}
				}
			case n := <-m.nodeWatcher.Queue:
				{
					fmt.Println(909090)
					fmt.Println(n.Event.Type)
					if n.Event.Type == mvccpb.DELETE {
						if val, ok := n.Model.(*node2.Node); ok {
							fmt.Printf("---> %v\n", val)
							pb, err := m.pipelineBind()
							if err != nil {
								blog.Error(err)
								continue
							}
							var bind bool
							var pipe string
							fmt.Println("end vvv aa")
							for pk, pv := range pb.Bindings {
								if pv == val.Name {
									bind = true
									pipe = pk
									break
								}
							}
							fmt.Println("end vvv")
							if bind {
								m.notBindPipelineCh <- &pipeline2.Pipeline{Name: pipe}
							}
						}
					}
					blog.Debug("monitor node queue \n", n)
				}
			case <-time.Tick(time.Second * 30):
				{
					m.queueAllPipelines()
				}
			}
		}
	}()
	return
}

func (m *Monitor) allPipelines() (res []*pipeline2.Pipeline, err error) {
	return dao.AllPipelines()
}

func (m *Monitor) allNodes() (res []*node2.Node, err error) {
	return dao.AllNodes()
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

func (m *Monitor) checkAllPipelines() (res []*pipeline2.Pipeline, err error) {
	res = []*pipeline2.Pipeline{}
	pipes, err := m.allPipelines()
	if err != nil {
		return
	}
	nodes, err := m.allNodes()
	if err != nil {
		return
	}
	pb, err := m.pipelineBind()
	if err != nil {
		return
	}
	mNodes := map[string]bool{}
	for _, v := range nodes {
		mNodes[v.Name] = true
	}
	for _, v := range pipes {
		if nodeName, ok := pb.Bindings[v.Name]; ok {
			if _, t := mNodes[nodeName]; !t {
				res = append(res, v)
			}
		} else {
			res = append(res, v)
		}
	}
	return
}

func (m *Monitor) queueAllPipelines() {
	pipes, err := m.checkAllPipelines()
	if err != nil {
		blog.Error("Check all pipelines error: ", err)
	}
	for _, v := range pipes {
		m.notBindPipelineCh <- v
	}
}
