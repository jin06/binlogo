package monitor

import (
	"context"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/jin06/binlogo/pkg/store/dao/dao_node"
	"github.com/jin06/binlogo/pkg/store/dao/dao_sche"
	"github.com/jin06/binlogo/pkg/watcher/node"
	"github.com/sirupsen/logrus"
	"time"
)

func (m *Monitor) monitorNode(ctx context.Context) error {
	watcher, err := node.New(dao_node.NodePrefix())
	if err != nil {
		return err
	}
	ch, err := watcher.WatchEtcdList(ctx)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			case e := <-ch:
				{
					if e.Event.Type == mvccpb.DELETE {
					}
				}
			case <-time.Tick(time.Second * 60):
				{
					if er := m.checkAllNodeBind(); er != nil {
						logrus.Error("Check all node bind error: ", er)
					}
				}
			}
		}
	}()
	return nil
}

func (m *Monitor) checkAllWorkNode() (err error) {
	regNodesMap, err := dao_node.AllRegisterNodesMap()
	if err != nil {
		return
	}
	statusMap, err := dao_node.StatusMap()
	if err != nil {
		return
	}
	for k, val := range statusMap {
		if _, ok := regNodesMap[k]; !ok {
			if val.Ready == true {
				_, err1 := dao_node.CreateOrUpdateStatus(k, dao_node.WithStatusReady(false))
				if err1 != nil {
					logrus.Error(err1)
				}
			}
		} else {
			if val.Ready == false {

			}
		}
	}
	return
}

func (m *Monitor) checkAllNodeBind() (err error) {
	nodes, err := dao_node.AllNodes()
	if err != nil {
		return
	}
	pb, err := dao_sche.GetPipelineBind()
	if err != nil {
		return
	}
	mNodes := map[string]bool{}
	for _, v := range nodes {
		mNodes[v.Name] = true
	}
	for k, v := range pb.Bindings {
		if _, ok := mNodes[v]; !ok {
			_, err = dao_sche.UpdatePipelineBind(k, "")
			if err != nil {
				return
			}
		}
	}
	return
}
