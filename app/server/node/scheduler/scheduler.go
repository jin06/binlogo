package scheduler

import (
	"context"
	"github.com/jin06/binlogo/pkg/store/model"
	"github.com/sirupsen/logrus"
	"sync"
)

type Scheduler struct {
	options *Options
	queue   *Queue
	lock    sync.Mutex
}

func (s *Scheduler) Run(ctx context.Context) {
	s._schedule(ctx)
	return
}

func (s *Scheduler) _schedule(ctx context.Context) {
	go func() {
		for {
			err := s.scheduleOne()
			if err != nil {
				logrus.Error(err)
			}
			select {
			case <-ctx.Done():
				{
					return
				}
			}
		}
	}()
	return
}

func (s *Scheduler) Stop(ctx context.Context) {
	return
}

func (s *Scheduler) scheduleOne() (err error) {
	p := s.queue.getOne()
	allNodes := []*model.Node{}
	a := algorithm{
		pipeline: p,
		allNodes: allNodes,
	}
	err = a.cal()
	if err != nil {
		logrus.Error(err)
	}
	return
}

func (s *Scheduler) watchPipeline (err error) {
	return
}

func New(opts ...Option) (s *Scheduler) {
	options := &Options{}
	s = &Scheduler{options: options}
	for _, v := range opts {
		v(options)
	}
	s.queue = newQueue()
	return
}

