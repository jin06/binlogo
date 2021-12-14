package node_status

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/watcher"
)

func withHandler(prefix string, name string) watcher.Handler {
	return func(e *clientv3.Event) (ev *watcher.Event, err error) {
		ev = &watcher.Event{}
		m := &node.Status{}
		ev.Event = e
		ev.Data = m
		if e.Type == mvccpb.DELETE {
			if name == "" {
				_, err = fmt.Sscanf(string(e.Kv.Key), prefix+"/%s", &m.NodeName)
				if err != nil {
					return
				}
			} else {
				m.NodeName = name
			}
		} else {
			err = json.Unmarshal(e.Kv.Value, m)
			if err != nil {
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

// New returns a new node status watcher
//func New(key string) (w *watcher.General, err error) {
//	w, err = watcher.NewGeneral(key)
//	if err != nil {
//		return
//	}
//	w.EventHandler = withHandler(key)
//	return
//}
