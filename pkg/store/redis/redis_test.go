package store_redis

import (
	"context"
	"testing"

	"github.com/jin06/binlogo/v2/configs"
)

func initTestRedis() {
	cfg := configs.Redis{
		Addr: "192.168.3.45",
		Port: 16379,
		DB:   0,
	}
	Init(context.Background(), cfg)
}

func TestRedis(t *testing.T) {
	initTestRedis()
	cmd := Default.client.Ping(context.Background())
	t.Log(cmd.String())
}

type Data struct {
	Name  string
	Value string
}

func TestUpdateField(t *testing.T) {
	initTestRedis()

	data := Data{}
	cmd := Default.client.HSet(context.Background(), "test2", &data)
	t.Log(cmd.Result())
	t.Log(cmd.Err())
}
