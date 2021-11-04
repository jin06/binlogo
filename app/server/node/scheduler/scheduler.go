package scheduler

import (
	"context"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/store/dao/dao_sche"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"sync"
	"time"
)

type Scheduler struct {
	options *Options
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

func (s *Scheduler) Run(ctx context.Context) (err error) {
	s.runLock.Lock()
	defer s.runLock.Unlock()
	if s.status == SCHEDULER_RUN {
		blog.Debug("Running, do nothing")
		return
	}
	cctx, cancel := context.WithCancel(ctx)
	err = s.watcher.run(cctx)
	if err != nil {
		return
	}
	s.cancel = cancel
	s._schedule(cctx)
	blog.Debug("scheduler.run")
	//s._monitor(ctx2)
	s.status = SCHEDULER_RUN
	return
}

func (s *Scheduler) _schedule(ctx context.Context) {
	go func() {
		for {
			select {
			case <-time.Tick(time.Second):
				{
				}
			case <-ctx.Done():
				{
					return
				}
			case p := <-s.watcher.notBindPipelineCh:
				{
					blog.Infof("%s not bind node, bind one ", p.Name)
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
	blog.Debugf("schedule one: %v", p)
	//p := s.queue.getOne()
	a := newAlgorithm(p)
	err = a.cal()
	if err != nil {
		return
	}
	blog.Debugf("best node for %s is %s \n", p.Name, a.bestNode.Name)
	_, err = dao_sche.UpdatePipelineBind(p.Name, a.bestNode.Name)
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
	s.runLock = sync.Mutex{}
	s.watcher, err = newWatcher()
	s.status = SCHEDULER_STOP
	if err != nil {
		blog.Error(err)
	}
	return
}
