package scheduler_binding

import (
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
