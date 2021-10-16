package scheduler

import (
	"context"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/store/dao"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Scheduler struct {
	options *Options
	queue   *Queue
	lock    sync.Mutex
	monitor *Monitor
	status  string
	runLock sync.Mutex
	cancel  context.CancelFunc
}

const (
	SCHEDULER_RUN  = "running"
	SCHEDULER_STOP = "stop"
)

func (s *Scheduler) Run(ctx context.Context) {
	s.runLock.Lock()
	defer s.runLock.Unlock()
	if s.status == SCHEDULER_RUN {
		blog.Debug("Running, do nothing")
		return
	}
	ctx2, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	s._schedule(ctx2)
	blog.Debug("scheduler.run")
	//s._monitor(ctx2)
	s.monitor.run(ctx2)
	s.status = SCHEDULER_RUN
	return
}

func (s *Scheduler) _schedule(ctx context.Context) {
	go func() {
		for {
			select {
			case <- time.Tick(time.Second):{
			}
			case <-ctx.Done():
				{
					return
				}
			//case p := <-s.queue.readyQueue:
			case p := <-s.monitor.notBindPipelineCh:
				{
					if err := s.scheduleOne(p); err != nil {
						blog.Error(err)
					}
				}
			}

		}
	}()
	return
}

func (s *Scheduler) Stop(ctx context.Context) {
	s.runLock.Lock()
	defer s.runLock.Unlock()
	if s.status == SCHEDULER_STOP {
		blog.Debug("Stopped, do nothing")
		return
	}
	s.cancel()
	return
}

func (s *Scheduler) scheduleOne(p *pipeline.Pipeline) (err error) {
	//blog.Debug("schedule one")
	//p := s.queue.getOne()
	a := newAlgorithm(p)
	err = a.cal()
	if err != nil {
		logrus.Error(err)
	}
	blog.Debugf("best node for %s is %s \n", p.Name, a.bestNode.Name)
	err = dao.UpdatePipelineBind(p.Name, a.bestNode.Name)
	if err != nil {
		blog.Error(err)
	}
	return
}

func (s *Scheduler) watchPipeline(err error) {
	return
}

func New(opts ...Option) (s *Scheduler) {
	var err error
	options := &Options{}
	s = &Scheduler{options: options}
	for _, v := range opts {
		v(options)
	}
	s.queue = newQueue()
	s.runLock = sync.Mutex{}
	s.monitor, err = newMonitor()
	s.status = SCHEDULER_STOP
	if err != nil {
		blog.Error(err)
	}
	return
}
