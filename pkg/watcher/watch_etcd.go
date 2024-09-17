package watcher

import (
	"context"
	"sync"

	"github.com/jin06/binlogo/v2/pkg/etcdclient"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// General a base watcher
// this will be use for most watchers
type General struct {
	options   *Options
	watchOnce sync.Once
	closeOnce sync.Once
	ch        chan *Event
}

// Handler function of handle event
type Handler func(*clientv3.Event) (*Event, error)

// New returns a new General watcher with handler
func New(opts ...Option) (w *General, err error) {
	w = &General{
		options: NewOptions(opts...),
		ch:      make(chan *Event, 1000),
	}
	return
}

// GetKey returns General key
func (w *General) GetKey() string {
	return w.options.Key
}

// WatchEtcd start watch etcd changes
func (w *General) WatchEtcd(ctx context.Context, opts ...clientv3.OpOption) (ch chan *Event, err error) {
	w.watchOnce.Do(func() {
		go func() {
			defer func() {
				w.Close()
			}()
			watchCh := etcdclient.Default().Watch(ctx, w.options.Key, opts...)
			//LOOP:
			for {
				select {
				case resp, ok := <-watchCh:
					{
						if !ok {
							return
							//break LOOP
						}
						if resp.Err() != nil {
							logrus.Error(resp.Err())
						}
						for _, val := range resp.Events {
							ev, err2 := w.options.Handler(val)
							if err2 != nil {
								logrus.Error("Event handle error: ", err2)
								continue
							}
							ch <- ev
						}
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
func (w *General) WatchEtcdList(ctx context.Context) (ch chan *Event, err error) {
	return w.WatchEtcd(ctx, clientv3.WithPrefix())
}

func (w *General) Close() error {
	w.closeOnce.Do(func() {
		close(w.ch)
	})
	return nil
}
