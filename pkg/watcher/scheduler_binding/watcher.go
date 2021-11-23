package scheduler_binding

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/jin06/binlogo/pkg/store/dao/dao_sche"
	"github.com/jin06/binlogo/pkg/store/model/scheduler"
	"github.com/jin06/binlogo/pkg/watcher"
	"github.com/sirupsen/logrus"
)

func withHandler() watcher.Handler {
	return func(e *clientv3.Event) (ev *watcher.Event, err error) {
		ev = &watcher.Event{}
		m := &scheduler.PipelineBind{}
		ev.Event = e
		ev.Data = m
		if e.Type == mvccpb.DELETE {
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

// New returns a new pipeline bind watcher
func New() (w *watcher.General, err error) {
	key := dao_sche.PipeBindPrefix()
	w, err = watcher.NewGeneral(key)
	if err != nil {
		return
	}
	w.EventHandler = withHandler()
	return
}
