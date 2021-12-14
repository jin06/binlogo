package str

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jin06/binlogo/pkg/watcher"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func withHandler(prefix string, name string) watcher.Handler {
	return func(e *clientv3.Event) (ev *watcher.Event, err error) {
		ev = &watcher.Event{}
		var m string
		ev.Event = e
		ev.Data = m
		if e.Type == mvccpb.DELETE {

		} else {
			err = json.Unmarshal(e.Kv.Value, &m)
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
