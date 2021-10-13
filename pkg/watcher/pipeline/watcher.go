package pipeline

import (
	"context"
	model2 "github.com/jin06/binlogo/pkg/store/model"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/jin06/binlogo/pkg/watcher"
)

type PipelineWatcher struct {
	*watcher.General
	modelFunc func() model2.Model
	Queue     chan *pipeline.Pipeline
}

func New(key string) (w *PipelineWatcher, err error) {
	w = &PipelineWatcher{}
	w.General, err = watcher.NewGeneral(key)
	w.Queue = make(chan *pipeline.Pipeline, 10000)
	w.modelFunc = func() model2.Model {
		return &pipeline.Pipeline{}
	}
	return
}

func (w *PipelineWatcher) putQueue(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			case m := <-w.General.Queue:
				{
					if val, ok := m.(*pipeline.Pipeline); ok {
						w.Queue <- val
					}
				}
			}
		}
	}()
}

func (w *PipelineWatcher) Watch(ctx context.Context) {
	w.General.WatchEtcd(ctx, w.modelFunc)
	w.putQueue(ctx)

}

func (w *PipelineWatcher) WatchList(ctx context.Context) {
	w.General.WatchEtcdList(ctx, w.modelFunc)
	w.putQueue(ctx)
}
