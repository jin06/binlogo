package scheduler_binding

import (
	"context"
	"errors"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/jin06/binlogo/pkg/store/model/scheduler"
	"github.com/jin06/binlogo/pkg/watcher"
	"github.com/sirupsen/logrus"
)

func withHandler(prefix string, name string) watcher.Handler {
	return func(e *clientv3.Event) (ev *watcher.Event, err error) {
		ev = &watcher.Event{}
		m := &scheduler.PipelineBind{}
		ev.Event = e
		ev.Data = m
		if e.Type != mvccpb.DELETE {
			err = m.Unmarshal(e.Kv.Value)
			if err != nil {
				logrus.Error(err)
				return
			}
		}
		return
	}
}

// Watch watch one change
func Watch(ctx context.Context, prefix string, name string) (ch chan *watcher.Event, err error) {
	if name == "" {
		err = errors.New("empty name")
		return
	}
	key := prefix + "/" + name
	w, err := watcher.New(key, withHandler(prefix, name))
	if err != nil {
		return
	}
	return w.WatchEtcd(ctx)
}

// WatchList watch list change
func WatchList(ctx context.Context, prefix string) (ch chan *watcher.Event, err error) {
	w, err := watcher.New(prefix, withHandler(prefix, ""))
	if err != nil {
		return
	}
	return w.WatchEtcdList(ctx)
}

// New returns a new pipeline bind watcher
//func New() (w *watcher.General, err error) {
//	key := dao_sche.PipeBindPrefix()
//	w, err = watcher.NewGeneral(key)
//	if err != nil {
//		return
//	}
//	w.EventHandler = withHandler()
//	return
//}
