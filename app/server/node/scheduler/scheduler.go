package scheduler

import (
	"context"
	"sync"
)

type Scheduler struct {
	options *Options
	queue *Queue
	lock sync.Mutex
}

func (s *Scheduler) Run(ctx context.Context) {
	return
}

func (s *Scheduler) Stop(ctx context.Context) {
	return
}

func New(opts ...Option) (s *Scheduler) {
	options := &Options{}
	s = &Scheduler{options: options}
	for _, v := range opts {
		v(options)
	}
	return
}


