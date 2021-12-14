package monitor

import (
	"context"

	"go.etcd.io/etcd/api/v3/mvccpb"
	"github.com/jin06/binlogo/pkg/store/dao/dao_sche"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/sirupsen/logrus"
)

func (m *Monitor) monitorRegister(ctx context.Context) (err error) {
	//ch, err := m.registerWatcher.WatchEtcdList(ctx)
	//if err != nil {
	//	return err
	//}
	ch, err := m.newNodeRegWatcherCh(ctx)
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
			case n := <-ch:
				{
					if n.Event.Type == mvccpb.DELETE {
						if val, ok := n.Data.(*node.Node); ok {
							pb, er := dao_sche.GetPipelineBind()
							if er != nil {
								logrus.Error(er)
								continue
							}
							var bind bool
							var pipe string
							for pk, pv := range pb.Bindings {
								if pv == val.Name {
									bind = true
									pipe = pk
									break
								}
							}
							if bind {
								_, err = dao_sche.UpdatePipelineBind(pipe, "")
								if err != nil {
									logrus.Error(err)
								}
							}
							//todo update node status
						}
					}
				}
			}
		}
	}()
	return nil
}
