package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jin06/binlogo/v2/pkg/watcher"

	"github.com/jin06/binlogo/v2/pkg/event"
	"github.com/jin06/binlogo/v2/pkg/store/dao/dao_sche"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
	"github.com/jin06/binlogo/v2/pkg/store/model/scheduler"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
)

// Scheduler schedule the pipeline.
// allocate the pipeline to the appropriate node for operation
type Scheduler struct {
	lock     sync.Mutex
	stopOnce sync.Once
	stopping chan struct{}
	Exit     bool
}

// New returns a new Scheduler
func New() (s *Scheduler) {
	s = &Scheduler{
		stopping: make(chan struct{}),
	}
	return
}

// Run start working
func (s *Scheduler) Run(ctx context.Context) (err error) {
	logrus.Info("scheduler run")
	defer func() {
		s.Exit = true
		logrus.Info("scheduler stop")
	}()
	s.schedule(ctx)
	return
}

func (s *Scheduler) schedule(ctx context.Context) {
	stx, cancel := context.WithCancel(ctx)
	defer cancel()
	key := fmt.Sprintf("%s/%s", dao_sche.SchedulerPrefix(), "pipeline_bind")
	w, err := watcher.New(watcher.WithHandler(watcher.WrapSchedulerBinding()), watcher.WithKey(key))
	defer w.Close()
	if err != nil {
		return
	}
	waCh, err := w.WatchEtcd(stx)
	if err != nil {
		return
	}
	schedulePipes(nil)
	ticker := time.NewTicker(time.Second * 60)
	defer ticker.Stop()
	for {
		select {
		case <-s.stopping:
			return
		case <-ctx.Done():
			return
		case ev, whOK := <-waCh:
			if !whOK {
				return
			}
			{
				if ev.Event.Type == mvccpb.PUT {
					if val, ok := ev.Data.(*scheduler.PipelineBind); ok {
						schedulePipes(val)
					}
				}
			}
		case <-ticker.C:
			schedulePipes(nil)
		}
	}
}

// Stop stop working
func (s *Scheduler) Stop() {
	s.stop()
	return
}

func (s *Scheduler) stop() {
	s.stopOnce.Do(func() {
		close(s.stopping)
	})
}

func scheduleOne(p *pipeline.Pipeline) (err error) {
	a := newAlgorithm(withAlgPipe(p))
	err = a.cal()
	if err != nil {
		return
	}
	_, err = dao_sche.UpdatePipelineBind(p.Name, a.bestNode.Name)
	if err != nil {
		return
	}
	return
}

// noScheduledPipelines find pipelines that not be scheduled
func noScheduledPipelines(pb *scheduler.PipelineBind) (pipes []*pipeline.Pipeline, err error) {
	if pb == nil {
		pb, err = dao_sche.GetPipelineBind()
		if err != nil {
			return
		}
	}
	pipes = []*pipeline.Pipeline{}
	for k, v := range pb.Bindings {
		if v == "" {
			pipes = append(pipes, &pipeline.Pipeline{Name: k})
		}
	}
	return
}

func schedulePipes(pb *scheduler.PipelineBind) {
	pipes, err := noScheduledPipelines(pb)
	if err != nil {
		return
	}
	for _, v := range pipes {
		er := scheduleOne(v)
		if er != nil {
			event.EventErrorPipeline(v.Name, "SCHEDULING: "+er.Error())
		} else {
			event.EventInfoPipeline(v.Name, "SCHEDULING SUCCESS")
		}
	}
}
