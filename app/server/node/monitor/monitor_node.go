package monitor

import (
	"context"
	"time"

	"github.com/jin06/binlogo/v2/pkg/store/dao"
	nodeM "github.com/jin06/binlogo/v2/pkg/store/model/node"
	"github.com/jin06/binlogo/v2/pkg/watcher"
	"github.com/sirupsen/logrus"
)

func (m *Monitor) monitorNode(ctx context.Context) (err error) {
	logrus.Info("monitor node run ")
	defer logrus.Info("monitor node stop")

	watcher, err := watcher.New()
	if err != nil {
		return
	}
	defer watcher.Close()
	regCh, err := watcher.WatchNode(ctx)
	if err != nil {
		return
	}

	if err := m.checkAllNode(ctx); err != nil {
		logrus.Error("Check all node bind error: ", err)
	}

	ticker := time.NewTicker(time.Second * 120)
	defer ticker.Stop()

	for {
		select {
		case <-m.closing:
			return
		case <-ctx.Done():
			{
				return
			}
		case e, ok := <-regCh:
			{
				if !ok {
					return
				}
				if err := handleEventRegNode(ctx, e); err != nil {
					logrus.Errorln(err)
				}
			}
		case <-ticker.C:
			{
				if err := m.checkAllNode(ctx); err != nil {
					logrus.Error("Check all node bind error: ", err)
				}
			}
		}
	}
}

func (m *Monitor) checkAllNode(ctx context.Context) (err error) {
	regNodesMap, err := dao.AllRegisterNodesMap(ctx)
	if err != nil {
		return
	}
	nodesMap, err := dao.AllNodesMap(ctx)
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
		_, err := dao.CreateOrUpdateStatus(
			ctx,
			k,
			nodeM.StatusConditions{
				nodeM.ConNetworkUnavailable: networkStat,
				nodeM.ConReady:              readyStat,
			},
		)
		if err != nil {
			logrus.Error(err)
		}
	}
	statusMap, err := dao.StatusMap(ctx)
	if err != nil {
		return
	}
	for k := range statusMap {
		if _, ok := nodesMap[k]; !ok {
			if err := dao.DeleteStatus(ctx, k); err != nil {
				logrus.Errorln(err)
			}
		}
	}
	return
}

func (m *Monitor) checkAllNodeBind(ctx context.Context) (err error) {
	nodes, err := dao.AllNodes(context.Background())
	if err != nil {
		return
	}
	pb, err := dao.GetPipelineBind(context.Background())
	if err != nil {
		return
	}
	mNodes := map[string]bool{}
	for _, v := range nodes {
		mNodes[v.Name] = true
	}
	for k, v := range pb.Bindings {
		if _, ok := mNodes[v]; !ok {
			if _, err = dao.UpdatePipelineBind(ctx, k, ""); err != nil {
				return
			}
		}
	}
	return
}

func handleEventRegNode(ctx context.Context, e *watcher.WatcherEvent) (err error) {
	if val, ok := e.Data.(*nodeM.Node); ok {
		if e.EventType.IsDelete() {
			_, err = dao.CreateOrUpdateStatus(
				ctx,
				val.Name,
				nodeM.StatusConditions{
					nodeM.ConReady:              false,
					nodeM.ConNetworkUnavailable: true,
				},
			)
			return
		}
		if e.EventType.IsUpdate() {
			_, err = dao.CreateOrUpdateStatus(
				ctx,
				val.Name,
				nodeM.StatusConditions{
					nodeM.ConReady:              true,
					nodeM.ConNetworkUnavailable: false,
				},
			)
			return
		}
	}
	return
}
