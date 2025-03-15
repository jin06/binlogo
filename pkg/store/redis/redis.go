package storeredis

import (
	"context"
	"fmt"

	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/pkg/store/model"
	"github.com/redis/go-redis/v9"
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
