package store_redis

import (
	"context"
	"fmt"

	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/pkg/store/model"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var Default *Redis

func Init(ctx context.Context, cfg configs.Redis) error {
	Default = NewRedis(ctx, cfg)
	return Default.client.Ping(ctx).Err()
}

func NewRedis(ctx context.Context, cfg configs.Redis) *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port),
		Password: cfg.Passwrod,
		DB:       cfg.DB,
	})
	r := &Redis{
		client: rdb,
	}
	r.prefix = fmt.Sprintf("/%s/%s", configs.APP, viper.GetString("cluster.name"))
	return r
}

type Redis struct {
	client *redis.Client
	prefix string
}

func (r *Redis) key(m model.Model) string {
	key := r.prefix + "/" + m.Key()
	return key
}

func (r *Redis) Create(ctx context.Context, m model.Model) (success bool, err error) {
	key := r.prefix + "/" + m.Key()
	val := m.Val()
	if err = r.client.Set(ctx, key, val, 0).Err(); err != nil {
		return
	}
	return
}

func (r *Redis) Update(ctx context.Context, m model.Model) (success bool, err error) {
	return
}

func (r *Redis) Delete(ctx context.Context, m model.Model) (success bool, err error) {
	return
}

func (r *Redis) Get(ctx context.Context, m model.Model) (success bool, err error) {
	cmd := r.client.Get(ctx, r.key(m))
	err = cmd.Err()
	return
}

func (r *Redis) List(ctx context.Context, list []model.Model) error {
	// cmd := r.client.HGetAll()
	return nil
}

func (r *Redis) getAllHashKeys(ctx context.Context) ([]string, error) {
	var hashKeys []string
	var cursor uint64
	rdb := r.client

	// 遍历所有键
	for {
		// 使用 SCAN 命令获取部分键
		keys, nextCursor, err := rdb.Scan(ctx, cursor, "*", 10).Result()
		if err != nil {
			return nil, err
		}

		// 检查每个键的类型，筛选出 hash 类型的键
		for _, key := range keys {
			keyType, err := rdb.Type(ctx, key).Result()
			if err != nil {
				return nil, err
			}
			if keyType == "hash" {
				hashKeys = append(hashKeys, key)
			}
		}

		// 如果 cursor 变为 0，表示遍历完成
		if nextCursor == 0 {
			break
		}

		// 更新游标
		cursor = nextCursor
	}

	return hashKeys, nil
}
