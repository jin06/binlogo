package storeredis

import (
	"context"
	"fmt"

	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/pkg/store/model"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var Default *Redis

func Init(ctx context.Context, cfg configs.Redis) error {
	Default = NewRedis(ctx, cfg)

	return Default.client.Ping(ctx).Err()
}

func NewRedis(ctx context.Context, cfg configs.Redis) *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	r := &Redis{
		client: rdb,
	}
	r.prefix = Prefix()
	return r
}

func GetClient() *redis.Client {
	return Default.GetClient()
}

type Redis struct {
	client *redis.Client
	prefix string
}

func (r *Redis) GetClient() *redis.Client {
	return r.client
}

func (r *Redis) key(m model.Model) string {
	key := r.prefix + "/" + m.Key()
	return key
}

func (r *Redis) getPrefix(m model.Model) string {
	return fmt.Sprintf("%s/%s", r.prefix, m.Key())
}

func (r *Redis) Create(ctx context.Context, m model.Model) (bool, error) {
	i, err := r.client.HSet(ctx, r.getPrefix(m), m).Result()
	ok := (i > 0)
	return ok, err
}

func (r *Redis) UpdateField(ctx context.Context, m model.Model, values map[string]any) (bool, error) {
	return r.client.HMSet(ctx, r.getPrefix(m), values).Result()
}

func (r *Redis) Update(ctx context.Context, m model.Model) (bool, error) {
	i, err := r.client.HSet(ctx, r.getPrefix(m), m).Result()
	ok := (i > 0)
	return ok, err
}

func (r *Redis) Delete(ctx context.Context, m model.Model) (err error) {
	return r.client.Del(ctx, r.getPrefix(m)).Err()
}

func (r *Redis) Get(ctx context.Context, m model.Model) (ok bool, err error) {
	cmd := r.client.HGetAll(ctx, r.key(m))
	err = cmd.Err()
	if err == redis.Nil {
		ok = false
		err = nil
		return
	}
	if err != nil {
		return
	}
	ok = true
	err = cmd.Scan(m)
	return
}

func (r *Redis) List(ctx context.Context, list []model.Model) error {
	// cmd := r.client.HGetAll()
	return nil
}

func (r *Redis) HashSet(ctx context.Context) {
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

func (r *Redis) GetField(ctx context.Context, key string, field string) (string, error) {
	str, err := r.client.HGet(ctx, key, field).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return str, nil
}

func AllDatas[T model.Model](ctx context.Context, key string, client redis.Client) (list []T, err error) {
	var cursor uint64
	for {
		keys, nextCursor, err := client.Scan(ctx, cursor, key, 10).Result()
		if err != nil {
			return nil, err
		}
		for _, key := range keys {
			keyType, err := client.Type(ctx, key).Result()
			if err != nil {
				logrus.Error(err)
				continue
			}
			if keyType == "hash" {
				cmd := client.HGetAll(ctx, key)
				if cmd.Err() != nil {
					return nil, cmd.Err()
				}
				var item T
				err := cmd.Scan(item)
				if err != nil {
					logrus.Error(err)
					continue
				}
				list = append(list, item)
			}
		}
		if nextCursor == 0 {
			break
		}
		cursor = nextCursor

	}

	return
}
