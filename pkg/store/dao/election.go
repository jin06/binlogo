package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
	store_redis "github.com/jin06/binlogo/v2/pkg/store/redis"
	"github.com/redis/go-redis/v9"
)

const (
	masterExpiration time.Duration = time.Second * 5
)

func masterKey() string {
	return fmt.Sprintf("/%s/%s", configs.Default.ClusterName, "master_lock")
}

func LeaseMasterLock(ctx context.Context) error {
	// store_redis.GetClient().Expire()
	return nil
}

func GetMasterLock(ctx context.Context) (string, error) {
	cmd := store_redis.GetClient().Get(ctx, masterKey())
	val, err := cmd.Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

func AcquireMasterLock(ctx context.Context, node *node.Node) error {
	cmd := store_redis.GetClient().SetNX(ctx, masterKey(), node.Name, masterExpiration)
	err := cmd.Err()
	if err != nil {
		return err
	}
	return err
}

func LeaderNode(ctx context.Context) (nodeName string, err error) {
	return myDao.LeaderNode(ctx)
}

func IsMaster(ctx context.Context, nodeName string) (bool, error) {
	value, err := store_redis.GetClient().Get(ctx, masterKey()).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return value == nodeName, nil
}

func MasterTTL(ctx context.Context) (time.Duration, error) {
	return store_redis.GetClient().TTL(ctx, masterKey()).Result()
}

func Lease(ctx context.Context) {
	// store_redis.GetClient().SetNX()
}

// AllElections returns all nodes in the election
func AllElections() (res []map[string]interface{}, err error) {
	return myDao.AllElections()
}
