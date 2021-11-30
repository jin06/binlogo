package scheduler

import (
	"context"
	"github.com/jin06/binlogo/pkg/event"
	"github.com/jin06/binlogo/pkg/store/dao/dao_sche"
	event2 "github.com/jin06/binlogo/pkg/store/model/event"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/sirupsen/logrus"
	"sync"
)

// Scheduler schedule the pipeline.
// allocate the pipeline to the appropriate node for operation
type Scheduler struct {
	lock    sync.Mutex
	watcher *Watcher
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
	s.watcher, err = newWatcher()
	if err != nil {
		return
	}
	err = s.watcher.run(myCtx)
	if err != nil {
		return
	}
	s.cancel = cancel
	go s._schedule(myCtx)
	logrus.Debug("scheduler.run")
	return
}

func (s *Scheduler) _schedule(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		case p := <-s.watcher.notBindPipelineCh:
			{
				logrus.Infof("%s not bind node, bind one ", p.Name)
				err := s.scheduleOne(p)
				if err != nil {
					logrus.Error(err)
				}
			}
		}
	}
}

// Stop stop working
func (s *Scheduler) Stop() {
	s.runLock.Lock()
	defer s.runLock.Unlock()
	if s.status == SCHEDULER_STOP {
		//logrus.Debug("Stopped, do nothing")
		return
	}
	s.cancel()
	return
}

func (s *Scheduler) scheduleOne(p *pipeline.Pipeline) (err error) {
	logrus.Debugf("schedule one: %v", p)
	//p := s.queue.getOne()
	a := newAlgorithm(withAlgPipe(p))
	defer func() {
		if err != nil {
			event.Event(event2.NewErrorPipeline(p.Name, "Scheduling error, "+err.Error()))
		} else {
			event.Event(event2.NewInfoPipeline(p.Name, p.Name+" is scheduled to node "+a.bestNode.Name))
		}
	}()
	err = a.cal()
	if err != nil {
		return
	}
	logrus.Debugf("best node for %s is %s \n", p.Name, a.bestNode.Name)
	_, err = dao_sche.UpdatePipelineBind(p.Name, a.bestNode.Name)
	if err != nil {
		logrus.Error(err)
	}
	return
}

func (s *Scheduler) watchPipeline(err error) {
	return
}

// New returns a new Scheduler
func New() (s *Scheduler) {
	s = &Scheduler{}
	s.runLock = sync.Mutex{}
	s.status = SCHEDULER_STOP
	return
}
