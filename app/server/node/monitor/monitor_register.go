package monitor

import (
	"context"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/store/dao/dao_sche"
	"github.com/jin06/binlogo/pkg/store/model/node"
)

func (m *Monitor) monitorRegister(ctx context.Context) {
	go func() {
		m.registerWatcher.WatchList(ctx)
		for {
			select {
				case <- ctx.Done():{
					return
				}
				case n := <- m.registerWatcher.Queue:{
					if n.Event.Type == mvccpb.DELETE {
						if val, ok := n.Data.(*node.Node); ok {
							pb, err := dao_sche.GetPipelineBind()
							if err != nil {
								blog.Error(err)
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
									blog.Error(err)
								}
							}
							//todo update node status
						}
					}
				}
			}
		}
	}()
}
