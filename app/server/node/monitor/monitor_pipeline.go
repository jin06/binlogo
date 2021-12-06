package monitor

import (
	"context"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
	"github.com/jin06/binlogo/pkg/store/dao/dao_sche"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/sirupsen/logrus"
	"time"
)

func (m *Monitor) monitorPipe(ctx context.Context) (resCtx context.Context, err error) {
	resCtx, cancel := context.WithCancel(ctx)
	defer func() {
		if err != nil {
			cancel()
		}
	}()
	ch, err := m.newPipeWatcherCh(ctx)
	if err != nil {
		return
	}
	m.checkAllPipelineBind()
	m.checkAllPipelineDelete()
	go func() {
		defer func() {
			cancel()
		}()
		for {
			select {
			case <-time.Tick(time.Second * 120):
				{
					m.checkAllPipelineBind()
					m.checkAllPipelineDelete()
				}
			case <-ctx.Done():
				{
					return
				}
			case n, ok := <-ch:
				{
					if !ok {
						return
					}
					if n.Event.Type == mvccpb.DELETE {
						if val, ok := n.Data.(*pipeline.Pipeline); ok {
							_, err1 := dao_sche.DeletePipelineBind(val.Name)
							if err1 != nil {
								logrus.Error("Delete pipeline bind failed: ", err1)
							}
						}
					}
					if n.Event.Type == mvccpb.PUT {
						if val, ok1 := n.Data.(*pipeline.Pipeline); ok1 {
							if val.IsDelete {
								m.deletePipeline(val)
								continue
							}
							if val.ExpectRun() {
								err1 := dao_sche.UpdatePipelineBindIfNotExist(val.Name, "")
								if err1 != nil {
									logrus.Error("Update pipeline bind failed ", err1)
								}
							} else {
								if _, err1 := dao_sche.DeletePipelineBind(val.Name); err1 != nil {
									logrus.Errorln(err1)
								}
							}
						}
					}
				}
			}
		}
	}()
	return
}

func (m *Monitor) checkAllPipelineBind() {
	var err error
	defer func() {
		if err != nil {
			logrus.Error("Check all pipeline bind error: ", err)
		}
	}()
	pipes, err := dao_pipe.AllPipelinesMap()
	if err != nil {
		return
	}
	pb, err := dao_sche.GetPipelineBind()
	if err != nil {
		return
	}

	for _, v := range pipes {
		if v.ExpectRun() {
			if _, ok := pb.Bindings[v.Name]; !ok {
				errUpdate := dao_sche.UpdatePipelineBindIfNotExist(v.Name, "")
				if errUpdate != nil {
					logrus.Error(errUpdate)
				}
			}
		} else {
			if _, ok := pb.Bindings[v.Name]; ok {
				_, errDel := dao_sche.DeletePipelineBind(v.Name)
				if errDel != nil {
					logrus.Errorln(errDel)
				}
			}
		}
	}

	pb, err = dao_sche.GetPipelineBind()
	if err != nil {
		return
	}

	for k := range pb.Bindings {
		if _, ok := pipes[k]; !ok {
			_, errDel := dao_sche.DeletePipelineBind(k)
			if errDel != nil {
				logrus.Errorln(errDel)
			}
		}
	}
	return
}

func (m *Monitor) checkAllPipelineDelete() {
	var err error
	defer func() {
		if err != nil {
			logrus.Error("Check all deleted pipelines error: ", err)
		}
	}()
	pipes, err := dao_pipe.AllPipelines()
	if err != nil {
		return
	}
	for _, v := range pipes {
		if !v.IsDelete {
			continue
		}
		errDelete := m.deletePipeline(v)
		if errDelete != nil {
			logrus.Errorln("Delete pipeline error, ", errDelete)
		}
	}
	return
}

func (m *Monitor) deletePipeline(p *pipeline.Pipeline) (err error) {
	if p.Status != pipeline.STATUS_STOP {
		var ok bool
		ok, err = dao_pipe.UpdatePipeline(p.Name, pipeline.WithPipeStatus(pipeline.STATUS_STOP))
		if err != nil || !ok {
			return
		}
	}
	ins, err := dao_pipe.GetInstance(p.Name)
	if err != nil {
		return
	}
	//if nodeName != "" {
	//	return
	//}
	if ins != nil {
		return
	}
	_, err = dao_pipe.DeletePosition(p.Name)
	if err != nil {
		return
	}
	_, err = dao_pipe.DeleteCompletePipeline(p.Name)
	if err != nil {
		return
	}
	return
}
