package watcher

import (
	"context"

	"github.com/jin06/binlogo/v2/pkg/store/model"
)

// Watcher a wrapper for etcd watcher
type Watcher interface {
	Watch(ctx context.Context, key string) (chan model.Model, error)
	WtachList(ctx context.Context, key string) (chan model.Model, error)
}
