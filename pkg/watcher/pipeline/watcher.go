package pipeline

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/jin06/binlogo/pkg/watcher"
	"github.com/sirupsen/logrus"
)

func withHandler(key string) watcher.Handler {
	return func(e *clientv3.Event) (ev *watcher.Event, err error) {
		ev = &watcher.Event{}
		m := &pipeline.Pipeline{}
		ev.Event = e
		ev.Data = m
		if e.Type == mvccpb.DELETE {
			_, err = fmt.Sscanf(string(e.Kv.Key), key +"/%s", &m.Name)
			if err != nil {
				logrus.Error(err)
			return
			}
		} else {
			err = m.Unmarshal(e.Kv.Value)
			if err != nil {
				logrus.Error(err)
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

