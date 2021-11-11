package monitor

import (
	"context"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
	"github.com/jin06/binlogo/pkg/store/dao/dao_sche"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/sirupsen/logrus"
)

func (m *Monitor) monitorPipe(ctx context.Context) (err error) {
	ch, err := m.pipeWatcher.WatchEtcdList(ctx)
	if err != nil {
		return
	}
	if er := m.checkAllPipelineBind(); er != nil {
		logrus.Error("Check all pipeline bind error: ", er)
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			case n := <-ch:
				{
					if n.Event.Type == mvccpb.DELETE {
						if val, ok := n.Data.(*pipeline.Pipeline); ok {
							_, err := dao_sche.DeletePipelineBind(val.Name)
							if err != nil {
								logrus.Error("Delete pipeline bind failed: ", err)
							}
						}
					}
					if n.Event.Type == mvccpb.PUT {
						if val, ok := n.Data.(*pipeline.Pipeline); ok {
							switch val.Status {
							case pipeline.STATUS_RUN:
								{
									err := dao_sche.UpdatePipelineBindIfNotExist(val.Name, "")
									if err != nil {
										logrus.Error("Update pipeline bind failed ", err)
									}
								}
							case pipeline.STATUS_STOP:
								{
									if _, err := dao_sche.DeletePipelineBind(val.Name); err != nil {
										logrus.Errorln(err)
									}
								}
							}
						}
					}
				}

				//case <-time.Tick(time.Second * 60):
				//	{
				//		if er := m.checkAllPipelineBind(); er != nil {
				//			logrus.Error("Check all pipeline bind error: ", er)
				//		}
				//	}
			}
		}
	}()
	return nil
}

func (m *Monitor) checkAllPipelineBind() (err error) {
	pipes, err := dao_pipe.AllPipelinesMap()
	if err != nil {
		return
	}
	pb, err := dao_sche.GetPipelineBind()
	if err != nil {
		return
	}

	for _, v := range pipes {
		switch v.Status {
		case pipeline.STATUS_STOP:
			{
				if _, ok := pb.Bindings[v.Name]; ok {
					_, errDel := dao_sche.DeletePipelineBind(v.Name)
					if errDel != nil {
						logrus.Errorln(errDel)
					}
				}
			}
		case pipeline.STATUS_RUN:
			{
				if _, ok := pb.Bindings[v.Name]; !ok {
					_, errUpdate := dao_sche.UpdatePipelineBind(v.Name, "")
					if errUpdate != nil {
						logrus.Error(errUpdate)
					}
				}
			}
		}
	}

	pb, err = dao_sche.GetPipelineBind()
	if err != nil {
		return
	}

	for k, _ := range pb.Bindings {
		if _, ok := pipes[k]; !ok {
			_, errDel := dao_sche.DeletePipelineBind(k)
			if errDel != nil {
				logrus.Errorln(errDel)
			}
		}
	}
	return
}
