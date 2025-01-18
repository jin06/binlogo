package watcher

import (
	"context"
	"fmt"

	"github.com/jin06/binlogo/v2/pkg/store/model"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
	storeredis "github.com/jin06/binlogo/v2/pkg/store/redis"
)

var def Watcher

func Default() Watcher {
	if def == nil {
		def = newRedisWatcher(storeredis.Default)
	}
	return def
}

// Watcher a wrapper for etcd watcher
type Watcher interface {
	// Watch(ctx context.Context, key string) (chan model.Model, error)
	// WtachList(ctx context.Context, key string) (chan model.Model, error)
	WatchPipelinBind(ctx context.Context) chan model.PipelineBind
	WatchNode(ctx context.Context) chan node.Node
}

func NewRedisWatcher(redis *storeredis.Redis) *RedisWatcher {
	return &RedisWatcher{
		redis: redis,
	}
}

func newRedisWatcher(redis *storeredis.Redis) *RedisWatcher {
	return &RedisWatcher{redis}
}

type RedisWatcher struct {
	redis *storeredis.Redis
}

func (watcher *RedisWatcher) WatchPipelinBind(ctx context.Context) chan model.PipelineBind {
	pubsub := watcher.redis.GetClient().Subscribe(ctx, storeredis.PipelineBindChan())
	ch := make(chan model.PipelineBind, 1000)
	go func() {
		msgch := pubsub.Channel()
		for {
			select {
			case msg := <-msgch:
				fmt.Println(msg)
			}
		}
	}()
	return ch
}

func (watcher *RedisWatcher) WatchNode(ctx context.Context) chan node.Node {
	pubsub := watcher.redis.GetClient().Subscribe(ctx, storeredis.PipelineBindChan())
	ch := make(chan node.Node, 1000)
	go func() {
		msgch := pubsub.Channel()
		for {
			select {
			case msg := <-msgch:
				fmt.Println(msg)
			}
		}
	}()
	return ch
}
