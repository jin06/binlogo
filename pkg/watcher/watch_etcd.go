package watcher

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/jin06/binlogo/pkg/etcd_client"
	"github.com/sirupsen/logrus"
)

type General struct {
	key          string
	EventHandler func(*clientv3.Event) (*Event, error)
}

type Handler func(*clientv3.Event) (*Event, error)

func NewGeneral(key string) (w *General, err error) {
	w = &General{
		key: key,
	}
	w.EventHandler = func(event *clientv3.Event) (*Event, error) {
		return nil, nil
	}
	return
}

func (w *General) GetKey() string {
	return w.key
}

func (w *General) WatchEtcd(ctx context.Context, opts ...clientv3.OpOption) (ch chan *Event, err error) {
	ch = make(chan *Event, 1000)
	defer func() {
		if err != nil {
			close(ch)
		}
	}()
	go func() {
		watchCh := etcd_client.Default().Watch(ctx, w.key, opts...)
	LOOP:
		for {
			select {
			case resp, ok := <-watchCh:
				{
					if !ok {
						break LOOP
					}
					if resp.Err() != nil {
						logrus.Error(resp.Err())
					}
					for _, val := range resp.Events {
						ev, err2 := w.EventHandler(val)
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
}

func (w *General) WatchEtcdList(ctx context.Context) (ch chan *Event, err error) {
	return w.WatchEtcd(ctx, clientv3.WithPrefix())
}
