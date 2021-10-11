package scheduler

import (
	"context"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
)

type Queue struct {
	readyQueue  chan *pipeline.Pipeline
	failedQueue chan *pipeline.Pipeline
}

func (q *Queue) getOne() (p *pipeline.Pipeline) {
	p = <-q.readyQueue
	return p
}

func (q *Queue) watchPipeline(ctx context.Context) {
}

func (q *Queue) putFailed(p *pipeline.Pipeline) {
	q.failedQueue <- p
	return
}

func newQueue() (q *Queue){
	q = &Queue{}
	q.readyQueue = make(chan *pipeline.Pipeline, 1000)
	q.failedQueue = make(chan *pipeline.Pipeline, 1000)
	return
}
