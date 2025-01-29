package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/jin06/binlogo/v2/pkg/watcher"

	"github.com/jin06/binlogo/v2/pkg/store/dao"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
	"github.com/sirupsen/logrus"
)

func (m *Monitor) monitorPipe(ctx context.Context) (err error) {
	logrus.Info("monitor pipeline run ")
	defer logrus.Info("monitor pipeline stop")
	ch := watcher.Default().WatchPipelinBind(ctx)

	m.checkAllPipelineBind(ctx)
	m.checkAllPipelineDelete(ctx)
	ticker := time.NewTicker(time.Second * 120)
	defer ticker.Stop()
	for {
		select {
		case <-m.closing:
			return
		case <-ticker.C:
			{
				m.checkAllPipelineBind(ctx)
				m.checkAllPipelineDelete(ctx)
			}
		case <-ctx.Done():
			return
		case bind, ok := <-ch:
			{
				if !ok {
					return
				}
				fmt.Println(bind)
				// if n.Event.Type == mvccpb.DELETE {
				// 	if val, ok := n.Data.(*pipeline.Pipeline); ok {
				// 		_, err1 := dao.DeletePipelineBind(ctx, val.Name)
				// 		if err1 != nil {
				// 			logrus.Error("Delete pipeline bind failed: ", err1)
				// 		}
				// 	}
				// }
				// if n.Event.Type == mvccpb.PUT {
				// 	if val, ok1 := n.Data.(*pipeline.Pipeline); ok1 {
				// 		if val.IsDelete {
				// 			m.deletePipeline(val)
				// 			continue
				// 		}
				// 		if val.ExpectRun() {
				// 			err1 := dao.UpdatePipelineBindIfNotExist(ctx, val.Name, "")
				// 			if err1 != nil {
				// 				logrus.Error("Update pipeline bind failed ", err1)
				// 			}
				// 		} else {
				// 			if _, err1 := dao.DeletePipelineBind(ctx, val.Name); err1 != nil {
				// 				logrus.Errorln(err1)
				// 			}
				// 		}
				// 	}
				// }
			}
		}
	}
}

func (m *Monitor) checkAllPipelineBind(ctx context.Context) {
	var err error
	defer func() {
		if err != nil {
			logrus.Error("Check all pipeline bind error: ", err)
		}
	}()
	pipes, err := dao.AllPipelinesMap(context.Background())
	if err != nil {
		return
	}
	pb, err := dao.GetPipelineBind(context.Background())
	if err != nil {
		return
	}

	for _, v := range pipes {
		if v.ExpectRun() {
			if _, ok := pb.Bindings[v.Name]; !ok {
				errUpdate := dao.UpdatePipelineBindIfNotExist(ctx, v.Name, "")
				if errUpdate != nil {
					logrus.Error(errUpdate)
				}
			}
		} else {
			if _, ok := pb.Bindings[v.Name]; ok {
				_, errDel := dao.DeletePipelineBind(ctx, v.Name)
				if errDel != nil {
					logrus.Errorln(errDel)
				}
			}
		}
	}

	pb, err = dao.GetPipelineBind(context.Background())
	if err != nil {
		return
	}

	for k := range pb.Bindings {
		if _, ok := pipes[k]; !ok {
			_, errDel := dao.DeletePipelineBind(ctx, k)
			if errDel != nil {
				logrus.Errorln(errDel)
			}
		}
	}
}

func (m *Monitor) checkAllPipelineDelete(ctx context.Context) {
	var err error
	defer func() {
		if err != nil {
			logrus.Error("Check all deleted pipelines error: ", err)
		}
	}()
	pipes, err := dao.AllPipelines(context.Background())
	if err != nil {
		return
	}
	for _, v := range pipes {
		if !v.IsDelete {
			continue
		}
		errDelete := m.deletePipeline(ctx, v)
		if errDelete != nil {
			logrus.Errorln("Delete pipeline error, ", errDelete)
		}
	}
}

func (m *Monitor) deletePipeline(ctx context.Context, p *pipeline.Pipeline) (err error) {
	if p.Status != pipeline.STATUS_STOP {
		if err := dao.UpdatePipeline(ctx, p.Name, pipeline.WithPipeStatus(pipeline.STATUS_STOP)); err != nil {
			return err
		}
	}
	ins, err := dao.GetInstance(context.Background(), p.Name)
	if err != nil || ins != nil {
		return
	}

	if err = dao.DeletePosition(ctx, p.Name); err != nil {
		return
	}
	err = dao.DeleteCompletePipeline(ctx, p.Name)
	return
}
