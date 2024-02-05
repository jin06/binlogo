package monitor

import (
	"context"
	"time"

	"github.com/jin06/binlogo/pkg/watcher"

	"github.com/jin06/binlogo/pkg/store/dao"
	"github.com/jin06/binlogo/pkg/store/dao/dao_node"
	"github.com/jin06/binlogo/pkg/store/dao/dao_sche"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
)

func (m *Monitor) monitorStatus(ctx context.Context) (err error) {
	logrus.Info("monitor status run ")
	defer logrus.Info("monitor status stop")
	stx, cancel := context.WithCancel(ctx)
	defer cancel()
	key := dao_node.StatusPrefix()
	w, err := watcher.New(watcher.WithKey(key), watcher.WithHandler(watcher.WrapNodeStatus(key, "")))
	defer w.Close()
	if err != nil {
		return
	}
	ch, err := w.WatchEtcdList(stx)
	if err != nil {
		return
	}
	if err = checkAllNodeStatus(); err != nil {
		return
	}

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		case <-ticker.C:
			{
				er := checkAllNodeStatus()
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
					if e.Event.Type == mvccpb.DELETE {
						err1 := removePipelineBindIfBindNode(val.NodeName)
						if err1 != nil {
							logrus.Errorln(err1)
						}
					}
					if e.Event.Type == mvccpb.PUT {
						if val.Ready == false {
							err1 := removePipelineBindIfBindNode(val.NodeName)
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

func checkAllNodeStatus() (err error) {
	mapping, err := dao_node.StatusMap()
	if err != nil {
		return
	}
	pb, err := dao_sche.GetPipelineBind()
	if err != nil {
		return
	}

	for k, v := range pb.Bindings {
		if val, ok := mapping[v]; !ok {
			err1 := dao.ClearOrDeleteBind(k)
			if err1 != nil {
				logrus.Error(err1)
			}
		} else {
			if val.Ready == false {
				err1 := dao.ClearOrDeleteBind(k)
				if err1 != nil {
					logrus.Errorln(err1)
				}
			}
		}
	}
	return
}

func removePipelineBindIfBindNode(nodeName string) (err error) {
	pb, err := dao_sche.GetPipelineBind()
	if err != nil {
		return
	}
	for k, v := range pb.Bindings {
		if v == nodeName {
			err = dao.ClearOrDeleteBind(k)
			break
		}
	}
	return
}
