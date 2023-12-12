package service

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var defaultStore TokenStore
var initDefaultStore sync.Once

func DefaultStore() TokenStore {
	initDefaultStore.Do(func() {
		duration := viper.GetDuration("auth.store.expiration")
		switch viper.GetString("auth.store") {
		case "redis":
			{
				client := redis.NewClient(&redis.Options{
					Addr:     viper.GetString("auth.store.redis.addr"),
					Username: viper.GetString("auth.store.redis.username"),
					Password: viper.GetString("auth.store.redis.password"),
					DB:       viper.GetInt("auth.store.redis.db"),
					PoolSize: 5,
				})
				defaultStore = &RedisStore{
					prefix:   "binlogo_token:",
					client:   client,
					duraiton: duration,
				}
			}
		case "memory":
			{
				defaultStore = &MemoryStore{
					tokens:   map[string]time.Time{},
					duration: duration,
				}
			}
		case "none":
			fallthrough
		default:
			defaultStore = &noneStore{}
		}
	})
	return defaultStore
}

type TokenStore interface {
	Set() (token string)
	Get(token string) bool
	Remove(token string)
}

type MemoryStore struct {
	mu       sync.Mutex
	tokens   map[string]time.Time
	duration time.Duration
}

func (m *MemoryStore) Set() (token string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	token = uuid.New().String()
	m.tokens[token] = time.Now()
	return token
}

func (m *MemoryStore) Get(token string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.tokens[token]; ok {
		if m.tokens[token].Add(m.duration).Before(time.Now()) {
			return false
		}
		m.tokens[token] = time.Now()
		return true
	}
	return false
}

func (m *MemoryStore) Remove(token string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.tokens, token)
}

type RedisStore struct {
	prefix   string
	duration time.Duration
	rdb      *redis.Client
}

func (r *RedisStore) Set() (token string) {
	token = uuid.New().String()
	r.rdb.SetEX(context.Background(), r.prefix+token, true, r.duration)
	return token
}

func (r *RedisStore) Get(token string) bool {
	err := r.rdb.Get(context.Background(), r.prefix+token).Err()
	if err != nil {
		return false
	}
	return true
}

func (r *RedisStore) Remove(token string) {
	r.rdb.Del(context.Background(), r.prefix+token)
}

type noneStore struct {
}

func (n *noneStore) Set() string {
	return uuid.New().String()
}

func (n *noneStore) Get(token string) bool {
	return true
}

func (n *noneStore) Remove(token string) {
	return
}
