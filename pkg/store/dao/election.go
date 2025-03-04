package dao

import (
	"context"

	"github.com/jin06/binlogo/v2/pkg/store/model/node"
)

func LeaseMasterLock(ctx context.Context) error {
	return myDao.LeaseMasterLock(ctx)
}

func GetMasterLock(ctx context.Context) (string, error) {
	return myDao.GetMasterLock(ctx)
}

func AcquireMasterLock(ctx context.Context, node *node.Node) error {
	return myDao.AcquireMasterLock(ctx, node)
}

func LeaderNode(ctx context.Context) (nodeName string, err error) {
	return myDao.LeaderNode(ctx)
}

// AllElections returns all nodes in the election
func AllElections() (res []map[string]interface{}, err error) {
	return myDao.AllElections()
}
