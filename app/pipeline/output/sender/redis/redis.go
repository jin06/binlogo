package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	message2 "github.com/jin06/binlogo/app/pipeline/message"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
)

// Reids send message to Redis
type Redis struct {
	Redis  *pipeline.Redis
	Client *redis.Client
}

// New returns a new Reids instance
func New(rs *pipeline.Redis) (r *Redis, err error) {
	r = &Redis{Redis: rs}
	err = r.init()
	return
}

func (r *Redis) init() (err error) {
	r.Client = redis.NewClient(&redis.Options{
		Addr:     r.Redis.Addr,
		Password: r.Redis.Password,
		DB:       r.Redis.DB,
	})
	err = r.ping()
	return
}

func (r *Redis) ping() (err error) {
	_, err = r.Client.Ping(context.Background()).Result()
	return
}

// Send loginc and control
func (r *Redis) Send(msg *message2.Message) (ok bool, err error) {
	body, err := msg.JsonContent()
	if err != nil {
		return
	}
	err = r.Client.RPush(context.Background(), r.Redis.List, body).Err()
	return
}
