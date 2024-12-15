package dao

import (
	"context"

	"github.com/jin06/binlogo/v2/pkg/etcdclient"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
	store_redis "github.com/jin06/binlogo/v2/pkg/store/redis"
)

// CapacityPrefix returns prefix key of etcd for capacity
func CapacityPrefix() string {
	return etcdclient.Prefix() + "/node/capacity"
}

// UpdateCapacity update capacity data to etcd
func UpdateCapacity(ctx context.Context, cap *node.Capacity) (bool, error) {
	return myDao.UpdateCapacity(ctx, cap)
}

// CapacityMap returns capacity data in map form
func CapacityMap(ctx context.Context) (mapping map[string]*node.Capacity, err error) {
	return myDao.CapacityMap(ctx)
}

// DeleteCapacity delete capacity in etcd
func DeleteCapacity(ctx context.Context, nodeName string) (bool, error) {
	return store_redis.Default.Delete(ctx, &node.Capacity{NodeName: nodeName})
}
