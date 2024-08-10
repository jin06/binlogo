package redis

import (
	"context"
	"fmt"

	"github.com/jin06/binlogo/v2/configs"
	"github.com/redis/go-redis/v9"
)

var Default Redis

type Redis struct {
	client *redis.Client
}

func NewRedis(cfg configs.Redis, ctx context.Context) *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port),
		Password: cfg.Passwrod,
		DB:       cfg.DB,
	})
	return &Redis{
		client: rdb,
	}
}

func DefaultRedis() {
	Default = *NewRedis(configs.Default.Store.Redis, context.Background())
}
