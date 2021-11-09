package scheduler

import (
	"context"
	"github.com/jin06/binlogo/pkg/store/dao/dao_sche"
	pipeline2 "github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/jin06/binlogo/pkg/store/model/scheduler"
	"github.com/jin06/binlogo/pkg/watcher"
	"github.com/sirupsen/logrus"
	"time"
)

type Watcher struct {
	pipelineWatcher   *watcher.General
	nodeWatcher       *watcher.General
	notBindPipelineCh chan *pipeline2.Pipeline
}

func newWatcher() (m *Watcher, err error) {
	m = &Watcher{
	}
	m.notBindPipelineCh = make(chan *pipeline2.Pipeline, 10000)
	return
}

func (m *Watcher) run(ctx context.Context) error {
	go func() {
		_ = m.putNotBindPipeToQueue(nil)
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			case <-time.Tick(60 * time.Second):
				{
					err := m.putNotBindPipeToQueue(nil)
					if err != nil {
						logrus.Error("Put not bind pipeline to queue error: ", err)
					}
				}
			}
			// todo watch pipeline bind
		}
	}()
	return nil
}

func (m *Watcher) putNotBindPipeToQueue(pb *scheduler.PipelineBind) (err error) {
	if pb == nil {
		pb, err = dao_sche.GetPipelineBind()
		if err != nil {
			return
		}
	}
	for k, v := range pb.Bindings {
		if v == "" {
			m.notBindPipelineCh <- &pipeline2.Pipeline{Name: k}
		}
	}
	return
}
