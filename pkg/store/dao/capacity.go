package dao

import (
	"context"

	"github.com/jin06/binlogo/v2/pkg/store/model/node"
)

// UpdateCapacity update capacity data to etcd
func UpdateCapacity(ctx context.Context, cap *node.Capacity) (bool, error) {
	return myDao.UpdateCapacity(ctx, cap)
}

// CapacityMap returns capacity data in map form
func CapacityMap(ctx context.Context) (mapping map[string]*node.Capacity, err error) {
	return myDao.CapacityMap(ctx)
}
