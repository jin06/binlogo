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
	m.nodeWatcher, err = node.New(prefix + "/register/nodes")
	return
}

func (m *Monitor) run(ctx context.Context) error {
	blog.Debug("monitor run")
	go func() {
		m.pipelineWatcher.WatchList(ctx)
		m.nodeWatcher.WatchList(ctx)
		if err := m.check(); err != nil {
			blog.Error(err)
		}
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			case n := <-m.pipelineWatcher.Queue:
				{
					if n.Event.Type == mvccpb.DELETE {
						if val, ok := n.Model.(*pipeline2.Pipeline); ok {
							_, err := dao.DeletePipelineBind(val.Name)
							if err != nil {
								blog.Error("Delete pipeline bind failed: ", err)
							}
						}
					}
					if n.Event.Type == mvccpb.PUT {
						if val, ok := n.Model.(*pipeline2.Pipeline); ok {
							err := dao.UpdatePipelineBindIfNotExist(val.Name, "")
							if err != nil {
								blog.Error("Update pipeline bind failed ", err)
							}
						}
					}
				}
			case n := <-m.nodeWatcher.Queue:
				{
					if n.Event.Type == mvccpb.DELETE {
						if val, ok := n.Model.(*node2.Node); ok {
							pb, err := m.pipelineBind()
							if err != nil {
								blog.Error(err)
								continue
							}
							var bind bool
							var pipe string
							for pk, pv := range pb.Bindings {
								if pv == val.Name {
									bind = true
									pipe = pk
									break
								}
							}
							if bind {
								//m.notBindPipelineCh <- &pipeline2.Pipeline{Name: pipe}
								_, err = dao.UpdatePipelineBind(pipe, "")
								if err != nil {
									blog.Error(err)
								}
							}
						}
					}
					blog.Debug("monitor node queue \n", n)
				}
			case <-time.Tick(time.Second * 60):
				{
					//m.queueAllPipelines()
					if err := m.check(); err != nil {
						blog.Error(err)
					}
				}
			case <-time.Tick(5 * time.Second):
				{
					err := m.putNotBindPipeToQueue()
					if err != nil {
						blog.Error("Put not bind pipeline to queue error: ", err)
					}
				}
			}
			// todo watch pipeline bind
		}
	}()
	return nil
}

func (m *Monitor) allPipelines() (res []*pipeline2.Pipeline, err error) {
	return dao.AllPipelines()
}

func (m *Monitor) allNodes() (res []*node2.Node, err error) {
	return dao.AllNodes()
}

func (m *Monitor) pipelineBind() (pb *scheduler.PipelineBind, err error) {
	return dao.GetPipelineBind()
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

func (m *Monitor) checkAllPipelineBind() (err error) {
	pipes, err := m.allPipelines()
	if err != nil {
		return
	}
	pb, err := m.pipelineBind()
	if err != nil {
		return
	}

	mapPipes := map[string]bool{}

	for _, v := range pipes {
		mapPipes[v.Name] = true
		if _, ok := pb.Bindings[v.Name]; !ok {
			_, err2 := dao.UpdatePipelineBind(v.Name, "")
			if err2 != nil {
				blog.Error(err2)
			}
		}
	}
	for k, _ := range pb.Bindings {
		if _, ok := mapPipes[k]; !ok {
			_, er := dao.DeletePipelineBind(k)
			if er != nil {
				blog.Error(er)
			}
		}
	}

	return
}

func (m *Monitor) checkAllNodeBind() (err error) {
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
	for k, v := range pb.Bindings {
		if _, ok := mNodes[v]; !ok {
			_, err = dao.UpdatePipelineBind(k, "")
			if err != nil {
				return
			}
		}
	}
	return
}

func (m *Monitor) check() (err error) {
	fmt.Println("---->")
	err = m.checkAllPipelineBind()
	if err != nil {
		return
	}
	err = m.checkAllNodeBind()
	fmt.Println("----> end")
	return
}

func (m *Monitor) putNotBindPipeToQueue() (err error) {
	pb, err := m.pipelineBind()
	if err != nil {
		return
	}
	for k, v := range pb.Bindings {
		if v == "" {
			fmt.Println(999)
			m.notBindPipelineCh <- &pipeline2.Pipeline{Name: k}
		}
	}
	return
}
