package dao

import (
	"context"

	"github.com/jin06/binlogo/v2/pkg/store/model/node"
)

// CreateOrUpdateStatus crate of update status in etcd
// create if not exist
func CreateOrUpdateStatus(ctx context.Context, nodeName string, conditions node.StatusConditions) (ok bool, err error) {
	return myDao.CreateOrUpdateStatus(ctx, nodeName, conditions)
}

// GetStatus get node status from etcd
func GetStatus(ctx context.Context, nodeName string) (s *node.Status, err error) {
	return myDao.GetStatus(ctx, nodeName)
}

// CreateStatusIfNotExist create status if not exist
func CreateStatusIfNotExist(ctx context.Context, nodeName string, condtions node.StatusConditions) (ok bool, err error) {
	return myDao.CreateOrUpdateStatus(ctx, nodeName, condtions)
}

// StatusMap returns all node status in map form
func StatusMap(ctx context.Context) (mapping map[string]*node.Status, err error) {
	return myDao.StatusMap(ctx)
}

// DeleteStatus delete node status in etcd
func DeleteStatus(ctx context.Context, name string) error {
	return myDao.DeleteStatus(ctx, name)
}
