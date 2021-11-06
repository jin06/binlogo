package watcher

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/jin06/binlogo/pkg/etcd_client"
	"github.com/sirupsen/logrus"
)

type General struct {
	Client  *clientv3.Client
	key   string
	Queue chan *clientv3.Event
}

func NewGeneral(key string) (w *General, err error) {
	w = &General{
		key:  key,
	}
	w.Queue = make(chan *clientv3.Event, 10000)
	w.Client, err = etcd_client.New()
	return
}

func (w *General) GetKey() string {
	return w.key
}

func (w *General) WatchEtcd(ctx context.Context, opts ...clientv3.OpOption) {
	watchCh := w.Client.Watch(ctx, w.key, opts...)
	go func() {
		for {
			select {
			case resp := <-watchCh:
				{
					if resp.Err() != nil {
						logrus.Error(resp.Err())
					}
					for _, val := range resp.Events {
						w.Queue <- val
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
}

func (w *General) WatchEtcdList(ctx context.Context) {
	w.WatchEtcd(ctx, clientv3.WithPrefix())
}
