package scheduler

import (
	"context"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/jin06/binlogo/pkg/store/dao/dao_sche"
	pipeline2 "github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/jin06/binlogo/pkg/store/model/scheduler"
	"github.com/jin06/binlogo/pkg/watcher"
	"github.com/jin06/binlogo/pkg/watcher/scheduler_binding"
	"github.com/sirupsen/logrus"
	"time"
)

type Watcher struct {
	pipelineWatcher   *watcher.General
	nodeWatcher       *watcher.General
	notBindPipelineCh chan *pipeline2.Pipeline
}

func newWatcher() (w *Watcher, err error) {
	w = &Watcher{}
	w.notBindPipelineCh = make(chan *pipeline2.Pipeline, 10000)
	return
}

func (w *Watcher) run(ctx context.Context) (err error) {
	wa, err := scheduler_binding.New()
	if err != nil {
		return
	}
	waCh, err := wa.WatchEtcd(ctx)
	if err != nil {
		return
	}
	// todo bug when disconnection
	go func() {
		_ = w.putNotBindPipeToQueue(nil)
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			case ev := <-waCh:
				{
					if ev.Event.Type == mvccpb.PUT {
						if val, ok := ev.Data.(*scheduler.PipelineBind); ok {
							errPut := w.putNotBindPipeToQueue(val)
							if errPut != nil {
								logrus.Error("Put not bind pipeline to queue error: ", errPut)
							}
						}
					}
				}
			case <-time.Tick(60 * time.Second):
				{
					err1 := w.putNotBindPipeToQueue(nil)
					if err1 != nil {
						logrus.Error(err1)
					}
				}
			}
			// todo watch pipeline bind
		}
	}()
	return nil
}

func (w *Watcher) putNotBindPipeToQueue(pb *scheduler.PipelineBind) (err error) {
	if pb == nil {
		pb, err = dao_sche.GetPipelineBind()
		if err != nil {
			return
		}
	}
	for k, v := range pb.Bindings {
		if v == "" {
			w.notBindPipelineCh <- &pipeline2.Pipeline{Name: k}
		}
	}
	return
}
