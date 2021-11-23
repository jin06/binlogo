package node_status

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/watcher"
)

func withHandler(key string) watcher.Handler {
	return func(e *clientv3.Event) (ev *watcher.Event, err error) {
		ev = &watcher.Event{}
		m := &node.Status{}
		ev.Event = e
		ev.Data = m
		if e.Type == mvccpb.DELETE {
			_, err = fmt.Sscanf(string(e.Kv.Key), key+"/%s", &m.NodeName)
			if err != nil {
				return
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

// New returns a new node status watcher
func New(key string) (w *watcher.General, err error) {
	w, err = watcher.NewGeneral(key)
	if err != nil {
		return
	}
	w.EventHandler = withHandler(key)
	return
}
