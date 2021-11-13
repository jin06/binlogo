package monitor

import (
	"context"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/jin06/binlogo/pkg/store/dao/dao_node"
	"github.com/jin06/binlogo/pkg/store/dao/dao_sche"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/watcher/node_status"
	"github.com/sirupsen/logrus"
	"time"
)

func (m *Monitor) monitorStatus(ctx context.Context) (err error) {
	watcher, err := node_status.New(dao_node.StatusPrefix())
	if err != nil {
		return
	}
	ch, err := watcher.WatchEtcdList(ctx)
	if err != nil {
		return
	}
	err = checkAllNodeStatus()
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
			case <-time.Tick(120 * time.Second):
				{
					err =checkAllNodeStatus()
				}
			case e := <-ch:
				{
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
			dao_sche.UpdatePipelineBind(k, "")
		} else {
			if val.Ready == false {
				dao_sche.UpdatePipelineBind(k, "")
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
			_, err = dao_sche.UpdatePipelineBind(k, "")
			break
		}
	}
	return
}
