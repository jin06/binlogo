package monitor

import (
	"context"
	"time"

	"github.com/jin06/binlogo/v2/pkg/watcher"

	"github.com/jin06/binlogo/v2/pkg/store/dao"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
	storeredis "github.com/jin06/binlogo/v2/pkg/store/redis"
	"github.com/sirupsen/logrus"
)

func (m *Monitor) monitorStatus(ctx context.Context) (err error) {
	logrus.Info("monitor status run ")
	defer logrus.Info("monitor status stop")
	stx, cancel := context.WithCancel(ctx)
	defer cancel()
	// key := dao.StatusPrefix()
	key := storeredis.StatusPrefix()
	w, err := watcher.New(watcher.WithKey(key))
	if err != nil {
		return
	}
	defer w.Close()
	ch, err := w.WatchEtcdList(stx)
	if err != nil {
		return
	}
	if err = checkAllNodeStatus(ctx); err != nil {
		return
	}

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-m.closing:
			return
		case <-ctx.Done():
			{
				return
			}
		case <-ticker.C:
			{
				er := checkAllNodeStatus(ctx)
				if er != nil {
					logrus.Errorln(er)
				}
			}
		case e, chOK := <-ch:
			{
				if !chOK {
					return
				}
				if val, ok := e.Data.(*node.Status); ok {
					if e.EventType == watcher.EventTypeDelete {
						err1 := removePipelineBindIfBindNode(ctx, val.NodeName)
						if err1 != nil {
							logrus.Errorln(err1)
						}
					}
					if e.EventType == watcher.EventTypeUpdate {
						if val.Ready == false {
							err1 := removePipelineBindIfBindNode(ctx, val.NodeName)
							if err1 != nil {
								logrus.Errorln(err1)
							}
						}
					}
				}
			}
		}
	}
}

func checkAllNodeStatus(ctx context.Context) (err error) {
	mapping, err := dao.StatusMap(ctx)
	if err != nil {
		return
	}
	pb, err := dao.GetPipelineBind(ctx)
	if err != nil {
		return
	}

	for k, v := range pb.Bindings {
		if val, ok := mapping[v]; !ok {
			err1 := dao.ClearOrDeleteBind(ctx, k)
			if err1 != nil {
				logrus.Error(err1)
			}
		} else {
			if val.Ready == false {
				err1 := dao.ClearOrDeleteBind(ctx, k)
				if err1 != nil {
					logrus.Errorln(err1)
				}
			}
		}
	}
	return
}

func removePipelineBindIfBindNode(ctx context.Context, name string) (err error) {
	pb, err := dao.GetPipelineBind(context.Background())
	if err != nil {
		return
	}
	for k, v := range pb.Bindings {
		if v == name {
			err = dao.ClearOrDeleteBind(ctx, k)
			break
		}
	}
	return
}
