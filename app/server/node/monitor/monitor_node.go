package monitor

import (
	"context"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/store/dao/dao_node"
	"github.com/jin06/binlogo/pkg/store/dao/dao_sche"
	"time"
)

func (m *Monitor) monitorNode(ctx context.Context) {
	go func() {
		m.nodeWatcher.WatchList(ctx)
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			case e := <-m.nodeWatcher.Queue:
				{
					if 	e.Event.Type == mvccpb.DELETE {

					}
				}
				case <- time.Tick(time.Second * 60): {
					if er := m.checkAllNodeBind(); er != nil {
						blog.Error("Check all node bind error: ", er)
					}
				}
			}
		}
	}()
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
