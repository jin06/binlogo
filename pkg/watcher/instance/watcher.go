package instance

import (
	"encoding/json"
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
		m := &pipeline.Instance{}
		ev.Event = e
		ev.Data = m
		_, err = fmt.Sscanf(string(e.Kv.Key), key +"/%s", &m.PipelineName)
		if err != nil {
			logrus.Error(err)
			return
		}
		if e.Type == mvccpb.DELETE {

		} else {
			err = json.Unmarshal(e.Kv.Value , &m.NodeName)
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

