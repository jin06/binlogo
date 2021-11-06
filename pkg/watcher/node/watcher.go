package node

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/watcher"
)

//type NodeWatcher struct {
//	Options *Options
//	*watcher.General
//}
//
//func New(key string) (w *NodeWatcher, err error) {
//	w = &NodeWatcher{}
//	w.General, err = watcher.NewGeneral(key)
//	return
//}
//
//func (w *NodeWatcher) _watch(ctx context.Context, option byte) (eventCh chan *watcher.Event, err error) {
//	var ch chan *clientv3.Event
//	if option == 1 {
//		ch, err = w.General.WatchEtcd(ctx)
//		if err != nil {
//			return
//		}
//	}
//	if option == 2 {
//		ch, err = w.General.WatchEtcdList(ctx)
//		if err != nil {
//			return
//		}
//	}
//	eventCh = make(chan *watcher.Event, 1000)
//	go func() {
//		for {
//			select {
//			case <-ctx.Done():
//				{
//					return
//				}
//			case e := <-ch:
//				{
//					val := &watcher.Event{
//					}
//					m := &node.Node{}
//					val.Event = e
//					val.Data = m
//					if e.Type == mvccpb.DELETE {
//						_, err1 := fmt.Sscanf(string(e.Kv.Key), w.General.GetKey()+"/%s", &m.Name)
//						fmt.Printf("%v", m)
//						if err1 != nil {
//							blog.Error(err1)
//							continue
//						}
//					} else {
//						err1 := m.Unmarshal(e.Kv.Value)
//						if err1 != nil {
//							blog.Error(err1)
//							continue
//						}
//					}
//					eventCh <- val
//				}
//			}
//		}
//	}()
//	return
//}
//
//func (w *NodeWatcher) Watch(ctx context.Context) (ch chan *watcher.Event, err error) {
//	return w._watch(ctx, 1)
//}
//
//func (w *NodeWatcher) WatchList(ctx context.Context) (ch chan *watcher.Event, err error) {
//	return w._watch(ctx, 2)
//}

func withHandler(key string) watcher.Handler {
	return func(e *clientv3.Event) (ev *watcher.Event, err error) {
		val := &watcher.Event{}
		m := &node.Node{}
		val.Event = e
		val.Data = m
		if e.Type == mvccpb.DELETE {
			_, err = fmt.Sscanf(string(e.Kv.Key), key+"/%s", &m.Name)
			if err != nil {
				return
			}
		} else {
			err = m.Unmarshal(e.Kv.Value)
			if err != nil {
				return
			}
		}
		return
	}
}

func New(key string) (w *watcher.General, err error) {
	w, err = watcher.NewGeneral(key)
	if err != nil {
		return
	}
	w.EventHandler = withHandler(key)
	return
}
