package watcher

import (
	"context"
	"github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
)

type General struct {
	Etcd  *etcd.ETCD
	key   string
	Queue chan *clientv3.Event
}

func NewGeneral(key string) (w *General, err error) {
	w = &General{
		key:  key,
		Etcd: etcd.E,
	}
	w.Queue = make(chan *clientv3.Event, 10000)
	return
}

func (w *General) GetKey() string {
	return w.key
}

func (w *General) WatchEtcd(ctx context.Context, opts ...clientv3.OpOption) {
	watchCh := w.Etcd.Client.Watch(ctx, w.key, opts...)
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
