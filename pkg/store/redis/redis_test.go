package store_redis

import (
	"context"
	"testing"

	"github.com/jin06/binlogo/v2/configs"
)

func TestRedis(t *testing.T) {
	cfg := configs.Redis{
		Addr: "192.168.3.45",
		Port: 16379,
		DB:   5,
	}

	r := NewRedis(cfg, context.Background())

	cmd := r.client.Ping(context.Background())

	t.Log(cmd.String())
}
