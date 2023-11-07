package scheduler

import (
	"context"
	"fmt"
	"github.com/jin06/binlogo/pkg/watcher"
	"sync"
	"time"

	"github.com/jin06/binlogo/pkg/event"
	"github.com/jin06/binlogo/pkg/store/dao/dao_sche"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/jin06/binlogo/pkg/store/model/scheduler"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
)

// Scheduler schedule the pipeline.
// allocate the pipeline to the appropriate node for operation
type Scheduler struct {
	lock    sync.Mutex
	status  string
	runLock sync.Mutex
	cancel  context.CancelFunc
}

const (
	SCHEDULER_RUN  = "run"
	SCHEDULER_STOP = "stop"
)

// Run start working
func (s *Scheduler) Run(ctx context.Context) (err error) {
	s.runLock.Lock()
	defer s.runLock.Unlock()
	if s.status == SCHEDULER_RUN {
		return
	}
	myCtx, cancel := context.WithCancel(ctx)
	defer func() {
		if err != nil {
			cancel()
		} else {
			s.status = SCHEDULER_RUN
		}
	}()
	s.cancel = cancel
	go s._schedule(myCtx)
	logrus.Debug("scheduler.run")
	return
}

func (s *Scheduler) _schedule(ctx context.Context) {
	defer func() {
		r := recover()
		if r != nil {
			logrus.Errorln("schedule error, ", r)
		}
		s.Stop()
	}()
	key := fmt.Sprintf("%s/%s", dao_sche.SchedulerPrefix(), "pipeline_bind")
	w, err := watcher.New(watcher.WithHandler(watcher.WrapSchedulerBinding()), watcher.WithKey(key))
	if err != nil {
		return
	}
	waCh, err := w.WatchEtcd(ctx)
	if err != nil {
		return
	}
	defer w.Close()
	schedulePipes(nil)
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
						schedulePipes(val)
					}
				}
			}
		case <-time.Tick(time.Second * 60):
			{
				schedulePipes(nil)
			}
		}
	}
}

// Stop stop working
func (s *Scheduler) Stop() {
	s.runLock.Lock()
	defer s.runLock.Unlock()
	if s.status == SCHEDULER_STOP {
		return
	}
	s.cancel()
	return
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

// New returns a new Scheduler
func New() (s *Scheduler) {
	s = &Scheduler{}
	s.runLock = sync.Mutex{}
	s.status = SCHEDULER_STOP
	s.cancel = func() {}
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

func (s *Scheduler) Status() string {
	return s.status
}
