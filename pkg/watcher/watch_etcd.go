package watcher

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/etcd_client"
	"github.com/sirupsen/logrus"
)

type General struct {
	key    string
	EventHandler func(*clientv3.Event) (*Event,error)
}

type Handler func(*clientv3.Event) (*Event,error)

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
	errCh := make(chan error)
	go func() {
		cli, err1 := etcd_client.New()
		if err1 != nil {
			errCh <- err1
		}
		defer func() {
			cli.Close()
		}()
		watchCh := cli.Watch(ctx, w.key, opts...)
		errCh <- nil
		for {
			select {
			case resp := <-watchCh:
				{
					if resp.Err() != nil {
						logrus.Error(resp.Err())
					}
					for _, val := range resp.Events {
						ev, err2 :=w.EventHandler(val)
						if err2 != nil {
							blog.Error("Event handle error: ", err2)
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
	err = <-errCh

	return
}

func (w *General) WatchEtcdList(ctx context.Context) (ch chan *Event, err error) {
	return w.WatchEtcd(ctx, clientv3.WithPrefix())
}
