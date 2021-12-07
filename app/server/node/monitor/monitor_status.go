package monitor

import (
	"context"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/jin06/binlogo/pkg/store/dao"
	"github.com/jin06/binlogo/pkg/store/dao/dao_node"
	"github.com/jin06/binlogo/pkg/store/dao/dao_sche"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/watcher/node_status"
	"github.com/sirupsen/logrus"
	"time"
)

func (m *Monitor) monitorStatus(ctx context.Context) (resCtx context.Context, err error) {
	resCtx, cancel := context.WithCancel(ctx)
	defer func() {
		if err != nil {
			cancel()
		}
	}()
	//watcher, err := node_status.New(dao_node.StatusPrefix())
	//if err != nil {
	//	return
	//}
	//ch, err := watcher.WatchEtcdList(ctx)
	ch, err := node_status.WatchList(ctx, dao_node.StatusPrefix())
	if err != nil {
		return
	}
	err = checkAllNodeStatus()
	if err != nil {
		return
	}
	go func() {
		defer func() {
			cancel()
		}()
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
