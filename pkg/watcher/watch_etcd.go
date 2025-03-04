package watcher

import (
	"context"
	"fmt"
	"sync"

	storeredis "github.com/jin06/binlogo/v2/pkg/store/redis"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// General a base watcher
// this will be use for most watchers
type General struct {
	options   *Options
	watchOnce sync.Once
	closeOnce sync.Once
	ch        chan *WatcherEvent
}

// Handler function of handle event
type Handler func(*clientv3.Event) (*Event, error)

// New returns a new General watcher with handler
func New(opts ...Option) (w *General, err error) {
	w = &General{
		options: NewOptions(opts...),
		ch:      make(chan *WatcherEvent, 1000),
	}
	return
}

// GetKey returns General key
func (w *General) GetKey() string {
	return w.options.Key
}

func (w *General) WatchNode(ctx context.Context) (chan *WatcherEvent, error) {
	w.options.Key = storeredis.NodeChan()
	return w.WatchEtcd(ctx)
}

// WatchEtcd start watch etcd changes
func (w *General) WatchEtcd(ctx context.Context) (ch chan *WatcherEvent, err error) {
	w.watchOnce.Do(func() {
		go func() {
			defer func() {
				w.Close()
			}()
			// watchCh := etcdclient.Default().Watch(ctx, w.options.Key)
			//LOOP:
			pubsub := storeredis.GetClient().Subscribe(ctx, w.options.Key)
			watchCh := pubsub.Channel()
			for {
				select {
				case resp, ok := <-watchCh:
					{
						if !ok {
							return
						}
						fmt.Println(resp)
						// if resp.Err() != nil {
						// 	logrus.Error(resp.Err())
						// }
						// for _, val := range resp.Events {
						// 	ev, err2 := w.options.Handler(val)
						// 	if err2 != nil {
						// 		logrus.Error("Event handle error: ", err2)
						// 		continue
						// 	}
						// 	ch <- ev
						// }
					}
				case <-ctx.Done():
					{
						return
					}
				}
			}

		}()
		return
	})
	return w.ch, nil
}

// WatchEtcdList start watch etcd changes for list
func (w *General) WatchEtcdList(ctx context.Context) (ch chan *WatcherEvent, err error) {
	return w.WatchEtcd(ctx)
}

func (w *General) Close() error {
	w.closeOnce.Do(func() {
		close(w.ch)
	})
	return nil
}
