package dao

import (
	"context"

	"github.com/jin06/binlogo/v2/pkg/store/model/node"
	// "github.com/jin06/binlogo/v2/pkg/store/storeredis"
)

func RegisterNode(ctx context.Context, n *node.Node) (bool, error) {
	return myDao.RegisterNode(ctx, n)
}

func LeaseNode(ctx context.Context, n *node.Node) error {
	return myDao.LeaseNode(ctx, n)
}

func RefreshNode(ctx context.Context, n *node.Node) (bool, error) {
	return myDao.RefreshNode(ctx, n)
}

func UpdateNodeIP(ctx context.Context, name string, ip string) (ok bool, err error) {
	return myDao.UpdateNodeIP(ctx, name, ip)
}

// GetNode get node data from etcd
func GetNode(ctx context.Context, name string) (*node.Node, error) {
	return myDao.GetNode(ctx, name)
}

// AllNodes return all node data from etcd
func AllNodes(ctx context.Context) (list []*node.Node, err error) {
	return myDao.AllNodes(ctx)
}

// AllNodesMap  returns all register nodes in map form
func AllNodesMap(ctx context.Context) (mapping map[string]*node.Node, err error) {
	list, err := AllNodes(ctx)
	if err != nil {
		return
	}
	mapping = map[string]*node.Node{}
	for _, v := range list {
		mapping[v.Name] = v
	}
	return
}

// AllRegisterNodesMap returns all register nodes from etcd in map form
func AllRegisterNodesMap(ctx context.Context) (mapping map[string]*node.Node, err error) {
	list, err := AllNodes(ctx)
	if err != nil {
		return
	}
	mapping = map[string]*node.Node{}
	for _, v := range list {
		mapping[v.Name] = v
	}
	return
}

// AllWorkNodesMap returns all register nodes from etcd in map form
func AllWorkNodesMap() (res map[string]*node.Node, err error) {
	res, err = AllNodesMap(context.Background())
	if err != nil {
		return
	}
	regMap, err := AllRegisterNodesMap(context.Background())
	if err != nil {
		return
	}
	for k := range res {
		if _, ok := regMap[k]; !ok {
			delete(res, k)
		}
	}
	return
}
