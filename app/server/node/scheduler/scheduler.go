package scheduler

import (
	"context"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/jin06/binlogo/pkg/event"
	"github.com/jin06/binlogo/pkg/store/dao/dao_sche"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/jin06/binlogo/pkg/store/model/scheduler"
	"github.com/jin06/binlogo/pkg/watcher/scheduler_binding"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
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
	defer s.Stop()
	wa, err := scheduler_binding.New()
	if err != nil {
		return
	}
	waCh, err := wa.WatchEtcd(ctx)
	if err != nil {
		return
	}
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
