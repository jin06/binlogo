package pipeline

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/jin06/binlogo/pkg/watcher"
)

//type PipelineWatcher struct {
//	*watcher.General
//	modelFunc func() model.Model
//}
//
//func New(key string) (w *PipelineWatcher, err error) {
//	w = &PipelineWatcher{}
//	w.General, err = watcher.NewGeneral(key)
//	return
//}
//
//func (w *PipelineWatcher) _watch(ctx context.Context, option byte)(eventCh chan *watcher.Event, err error) {
//	if option == 1 {
//		w.General.WatchEtcd(ctx)
//	}
//	if option == 2 {
//		w.General.WatchEtcdList(ctx)
//	}
//	go func() {
//		for {
//			select {
//			case <-ctx.Done():
//				{
//					return
//				}
//			case e := <-w.General.Queue:
//				{
//					val := &watcher.Event{
//					}
//					m := &pipeline.Pipeline{}
//					val.Event = e
//					val.Data = m
//					if e.Type == mvccpb.DELETE {
//						_, err := fmt.Sscanf(string(e.Kv.Key), w.General.GetKey()+"/%s", &m.Name)
//						if err != nil {
//							blog.Error(err)
//							continue
//						}
//					} else {
//						err := m.Unmarshal(e.Kv.Value)
//						if err != nil {
//							blog.Error(err)
//							continue
//						}
//					}
//					w.Queue <- val
//				}
//			}
//		}
//	}()
//}
//
//func (w *PipelineWatcher) Watch(ctx context.Context) {
//	w._watch(ctx, 1)
//}
//
//func (w *PipelineWatcher) WatchList(ctx context.Context) {
//	w._watch(ctx, 2)
//}

func withHandler(key string) watcher.Handler {
	return func(e *clientv3.Event) (ev *watcher.Event, err error) {
		val := &watcher.Event{
		}
		m := &pipeline.Pipeline{}
		val.Event = e
		val.Data = m
		if e.Type == mvccpb.DELETE {
			_, err = fmt.Sscanf(string(e.Kv.Key), key +"/%s", &m.Name)
			if err != nil {
				blog.Error(err)
			return
			}
		} else {
			err = m.Unmarshal(e.Kv.Value)
			if err != nil {
				blog.Error(err)
				return
			}
		}
		return
	}
}

func New(key string) (w *watcher.General, err error) {
	w, err  = watcher.NewGeneral(key)
	if err != nil {
		return
	}
	w.EventHandler = withHandler(key)
	return
}

