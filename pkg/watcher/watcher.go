package watcher

import (
	"github.com/jin06/binlogo/pkg/store/model"
)

// Watcher a wrapper for etcd watcher
type Watcher interface {
	Watch() chan model.Model
}
