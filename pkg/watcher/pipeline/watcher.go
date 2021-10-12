package pipeline

import (
	"context"
	model2 "github.com/jin06/binlogo/pkg/store/model"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/jin06/binlogo/pkg/watcher"
)

type PipelineWatcher struct {
	*watcher.General
}

func New(key string) (w *PipelineWatcher, err error) {
	w = &PipelineWatcher{}
	w.General, err = watcher.NewGeneral(key)
	return
}

func (w *PipelineWatcher) Watch(ctx context.Context) {
	w.General.WatchEtcd(ctx, func() model2.Model {
		return &pipeline.Pipeline{}
	})

}
