package dao

import (
	"context"

	"github.com/jin06/binlogo/v2/pkg/etcdclient"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
	store_redis "github.com/jin06/binlogo/v2/pkg/store/redis"
)

// AllocatablePrefix returns prefix key of etcd for allocatable
func AllocatablePrefix() string {
	return etcdclient.Prefix() + "/node/allocatable"
}

// UpdateAllocatable update allocatable data to etcd
func UpdateAllocatable(ctx context.Context, al *node.Allocatable) (ok bool, err error) {
	return store_redis.Default.Update(ctx, al)
}

// DeleteAllocatable delete allocatable data in etcd
func DeleteAllocatable(ctx context.Context, nodeName string) (ok bool, err error) {
	return store_redis.Default.Delete(ctx, &node.Allocatable{NodeName: nodeName})
}
