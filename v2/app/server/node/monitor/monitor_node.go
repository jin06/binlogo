package monitor

import (
	"context"
	"time"

	"github.com/jin06/binlogo/v2/pkg/store/dao/dao_node"
	"github.com/jin06/binlogo/v2/pkg/store/dao/dao_sche"
	nodeModel "github.com/jin06/binlogo/v2/pkg/store/model/node"
	"github.com/jin06/binlogo/v2/pkg/watcher"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
)

func (m *Monitor) monitorNode(ctx context.Context) (err error) {
	logrus.Info("monitor node run ")
	defer logrus.Info("monitor node stop")
	stx, cancel := context.WithCancel(ctx)
	defer cancel()
	key := dao_node.NodePrefix()

	nodeWatcher, err := watcher.New(watcher.WithKey(key), watcher.WithHandler(watcher.WrapNode(key, "")))
	defer nodeWatcher.Close()
	if err != nil {
		return
	}
	ch, err := nodeWatcher.WatchEtcdList(stx)
	if err != nil {
		return
	}

	regKey := dao_node.NodeRegisterPrefix()
	regWatcher, err := watcher.New(watcher.WithKey(regKey), watcher.WithHandler(watcher.WrapNode(regKey, "")))
	defer regWatcher.Close()
	if err != nil {
		return
	}

	regCh, err := regWatcher.WatchEtcdList(ctx)
	if err != nil {
		return
	}

	if er := m.checkAllNode(); er != nil {
		logrus.Error("Check all node bind error: ", er)
	}

	ticker := time.NewTicker(time.Second * 120)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		case e, chOK := <-ch:
			{
				if !chOK {
					return
				}
				if e.Event.Type == mvccpb.DELETE {
				}
			}
		case e, ok := <-regCh:
			{
				if !ok {
					return
				}
				er := handleEventRegNode(e)
				if er != nil {
					logrus.Errorln(er)
				}
			}
		case <-ticker.C:
			{
				if er := m.checkAllNode(); er != nil {
					logrus.Error("Check all node bind error: ", er)
				}
			}
		}
	}
}

func (m *Monitor) checkAllNode() (err error) {
	regNodesMap, err := dao_node.AllRegisterNodesMap()
	if err != nil {
		return
	}
	nodesMap, err := dao_node.AllNodesMap()
	if err != nil {
		return
	}
	for k := range nodesMap {
		_, ok := regNodesMap[k]
		readyStat := false
		networkStat := true
		if ok {
			readyStat = true
			networkStat = false
		}
		_, err1 := dao_node.CreateOrUpdateStatus(k, nodeModel.WithReady(readyStat), nodeModel.WithNetworkUnavailable(networkStat))
		if err1 != nil {
			logrus.Error(err1)
		}
	}
	statusMap, err := dao_node.StatusMap()
	if err != nil {
		return
	}
	for k := range statusMap {
		if _, ok := nodesMap[k]; !ok {
			_, err1 := dao_node.DeleteStatus(k)
			if err1 != nil {
				logrus.Errorln(err1)
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

func handleEventRegNode(e *watcher.Event) (err error) {
	if val, ok := e.Data.(*nodeModel.Node); ok {
		if e.Event.Type == mvccpb.DELETE {
			_, err = dao_node.CreateOrUpdateStatus(val.Name, nodeModel.WithReady(false), nodeModel.WithNetworkUnavailable(true))
			return
		}
		if e.Event.Type == mvccpb.PUT {
			_, err = dao_node.CreateOrUpdateStatus(val.Name, nodeModel.WithReady(true), nodeModel.WithNetworkUnavailable(false))
			return
		}
	}
	return
}
