package watcher

import (
	"context"
	"fmt"
	"github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/jin06/binlogo/pkg/store/model"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/sirupsen/logrus"
)

type Watcher interface {
	Watch() chan model.Model
}

type General struct {
	Etcd  *etcd.ETCD
	queue chan *pipeline.Pipeline
	key   string
}

func New(key string) (w *General, err error) {
	w = &General{
		key: key,
	}
	err = w.Init()
	return
}

func (w *General) Init() (err error) {
	w.Etcd = etcd.E
	return
}

func (w *General) Watch(ctx context.Context) {
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
						fmt.Printf("%v \n", val)
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
