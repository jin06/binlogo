package monitor

import (
	"context"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
	"github.com/jin06/binlogo/pkg/store/dao/dao_sche"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"time"
)

func (m *Monitor) monitorPipe(ctx context.Context) (err error){
	ch, err := m.pipeWatcher.WatchEtcdList(ctx)
	if err != nil {
		return
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			case n := <- ch:
				{
					if n.Event.Type == mvccpb.DELETE {
						if val, ok := n.Data.(*pipeline.Pipeline); ok {
							_, err := dao_sche.DeletePipelineBind(val.Name)
							if err != nil {
								blog.Error("Delete pipeline bind failed: ", err)
							}
						}
					}
					if n.Event.Type == mvccpb.PUT {
						if val, ok := n.Data.(*pipeline.Pipeline); ok {
							err := dao_sche.UpdatePipelineBindIfNotExist(val.Name, "")
							if err != nil {
								blog.Error("Update pipeline bind failed ", err)
							}
						}
					}
				}

			case <-time.Tick(time.Second * 60):
				{
					if er := m.checkAllPipelineBind(); er != nil {
						blog.Error("Check all pipeline bind error: ", er)
					}
				}
			}
		}
	}()
	return nil
}

func (m *Monitor) checkAllPipelineBind() (err error) {
	pipes, err := dao_pipe.AllPipelines()
	if err != nil {
		return
	}
	pb, err := dao_sche.GetPipelineBind()
	if err != nil {
		return
	}

	mapPipes := map[string]bool{}

	for _, v := range pipes {
		mapPipes[v.Name] = true
		if _, ok := pb.Bindings[v.Name]; !ok {
			_, err2 := dao_sche.UpdatePipelineBind(v.Name, "")
			if err2 != nil {
				blog.Error(err2)
			}
		}
	}
	for k, _ := range pb.Bindings {
		if _, ok := mapPipes[k]; !ok {
			_, er := dao_sche.DeletePipelineBind(k)
			if er != nil {
				blog.Error(er)
			}
		}
	}

	return
}
