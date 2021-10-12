package scheduler_binding

import (
	"context"
	model2 "github.com/jin06/binlogo/pkg/store/model"
	"github.com/jin06/binlogo/pkg/store/model/scheduler"
	"github.com/jin06/binlogo/pkg/watcher"
)

func New(key string) (w *BindingWatcher, err error) {
	w = &BindingWatcher{}
	w.General, err = watcher.NewGeneral(key)
	return
}

type BindingWatcher struct {
	*watcher.General
}

func (w *BindingWatcher) Watch(ctx context.Context) {
	w.General.WatchEtcd(ctx, func() model2.Model {
		return &scheduler.PipelineBind{}
	})
}
