package node

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/watcher"
)

type NodeWatcher struct {
	Options *Options
	*watcher.General
	Queue chan *watcher.Event
}

func New(key string) (w *NodeWatcher, err error) {
	w = &NodeWatcher{}
	w.General, err = watcher.NewGeneral(key)
	w.Queue = make(chan *watcher.Event, 1000)
	return
}

func (w *NodeWatcher) _watch(ctx context.Context, option byte) {
	if option == 1 {
		w.General.WatchEtcd(ctx)
	}
	if option == 2 {
		w.General.WatchEtcdList(ctx)
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			case e := <-w.General.Queue:
				{
					val := &watcher.Event{
					}
					m := &node.Node{}
					val.Event = e
					val.Data = m
					if e.Type == mvccpb.DELETE {
						_, err := fmt.Sscanf(string(e.Kv.Key), w.General.GetKey()+"/%s", &m.Name)
						fmt.Printf("%v", m)
						if err != nil {
							blog.Error(err)
							continue
						}
					} else {
						err := m.Unmarshal(e.Kv.Value)
						if err != nil {
							blog.Error(err)
							continue
						}
					}
					w.Queue <- val
				}
			}
		}
	}()
}

func (w *NodeWatcher) Watch(ctx context.Context) {
	w._watch(ctx, 1)
}

func (w *NodeWatcher) WatchList(ctx context.Context) {
	w._watch(ctx, 2)
}
