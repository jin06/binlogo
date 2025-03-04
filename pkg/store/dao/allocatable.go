package dao

import (
	"context"

	"github.com/jin06/binlogo/v2/pkg/store/model/node"
)

// UpdateAllocatable update allocatable data to etcd
func UpdateAllocatable(ctx context.Context, al *node.Allocatable) (ok bool, err error) {
	return myDao.UpdateAllocatable(ctx, al)
}
