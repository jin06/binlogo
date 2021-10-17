package watcher

import (
	"context"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/jin06/binlogo/pkg/store/model"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
)

type General struct {
	Etcd  *etcd.ETCD
	key   string
	Queue chan model.Model
	Queue2 chan *clientv3.Event
}

func NewGeneral(key string) (w *General, err error) {
	w = &General{
		key:   key,
		Queue: make(chan model.Model, 10000),
		Etcd:  etcd.E,
	}
	return
}

func (w *General) GetKey() string{
	return w.key
}

func (w *General) WatchEtcd(ctx context.Context, object func() model.Model) {
	watchCh := w.Etcd.Client.Watch(ctx, w.key)
	go func() {
		for {
			select {
			case resp := <-watchCh:
				{
					if resp.Err() != nil {
						logrus.Error(resp.Err())
					}
					for _, val := range resp.Events {
						obj := object()
						err := obj.Unmarshal(val.Kv.Value)
						if err != nil {
							blog.Error(err)
						}
						w.Queue <- obj
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

func (w *General) WatchEtcdList(ctx context.Context, object func() model.Model) {
	//fmt.Println(w.key)
	watchCh := w.Etcd.Client.Watch(ctx, w.key, clientv3.WithPrefix())
	go func() {
		for {
			select {
			case resp := <- watchCh:
				{
					if resp.Err() != nil {
						logrus.Error(resp.Err())
					}
					for _, val := range resp.Events {
						obj := object()
						err := obj.Unmarshal(val.Kv.Value)
						if err != nil {
							blog.Error(err)
						}
						//logrus.Debug("watch list object : ", obj)
						w.Queue <- obj
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

func (w *General) WatchEtcd2(ctx context.Context, opts ...clientv3.OpOption) chan *clientv3.Event{
	watchCh := w.Etcd.Client.Watch(ctx, w.key, opts...)
	eventCh := make(chan *clientv3.Event, 10000)
	w.Queue2 = eventCh
	go func() {
		for {
			select {
			case resp := <- watchCh:
				{
					if resp.Err() != nil {
						logrus.Error(resp.Err())
					}
					for _, val := range resp.Events {
						w.Queue2 <- val
					}
				}
			case <-ctx.Done():
				{
					return
				}
			}
		}
	}()
	return eventCh
}

func (w *General) WatchEtcdList2(ctx context.Context) chan *clientv3.Event {
	return w.WatchEtcd2(ctx, clientv3.WithPrefix())
}

