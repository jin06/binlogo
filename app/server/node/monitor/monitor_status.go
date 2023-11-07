package monitor

import (
	"context"
	"github.com/jin06/binlogo/pkg/watcher"
	"time"

	"github.com/jin06/binlogo/pkg/store/dao"
	"github.com/jin06/binlogo/pkg/store/dao/dao_node"
	"github.com/jin06/binlogo/pkg/store/dao/dao_sche"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
)

func (m *Monitor) monitorStatus(ctx context.Context) (resCtx context.Context, err error) {
	resCtx, cancel := context.WithCancel(ctx)
	defer func() {
		if err != nil {
			cancel()
		}
	}()
	key := dao_node.StatusPrefix()
	w, err := watcher.New(watcher.WithKey(key), watcher.WithHandler(watcher.WrapNodeStatus(key, "")))
	if err != nil {
		return
	}
	ch, err := w.WatchEtcdList(ctx)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			w.Close()
		}
	}()
	err = checkAllNodeStatus()
	if err != nil {
		return
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorln("monitor status panic", r)
			}
			cancel()
		}()
		defer w.Close()
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			case <-time.Tick(60 * time.Second):
				{
					er := checkAllNodeStatus()
					if er != nil {
						logrus.Errorln(er)
					}
				}
			case e, ok1 := <-ch:
				{
					if !ok1 {
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
	}()
	return
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
