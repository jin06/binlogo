package monitor

import (
	"context"
	"time"

	"github.com/jin06/binlogo/v2/pkg/watcher"

	"github.com/jin06/binlogo/v2/pkg/store/dao/dao_pipe"
	"github.com/jin06/binlogo/v2/pkg/store/dao/dao_sche"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
)

func (m *Monitor) monitorPipe(ctx context.Context) (err error) {
	logrus.Info("monitor pipeline run ")
	defer logrus.Info("monitor pipeline stop")
	stx, cancel := context.WithCancel(ctx)
	defer cancel()
	key := dao_pipe.PipelinePrefix()
	w, err := watcher.New(watcher.WithKey(key), watcher.WithHandler(watcher.WrapPipeline(key, "")))
	defer w.Close()
	if err != nil {
		return
	}
	//ch, err := m.newPipeWatcherCh(ctx)
	ch, err := w.WatchEtcdList(stx)
	if err != nil {
		return
	}

	m.checkAllPipelineBind()
	m.checkAllPipelineDelete()
	ticker := time.NewTicker(time.Second * 120)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			{
				m.checkAllPipelineBind()
				m.checkAllPipelineDelete()
			}
		case <-ctx.Done():
			return
		case n, chOK := <-ch:
			{
				if !chOK {
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
