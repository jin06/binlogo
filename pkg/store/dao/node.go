package dao

import (
	"context"
	"time"

	"github.com/jin06/binlogo/v2/pkg/etcdclient"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
	store_redis "github.com/jin06/binlogo/v2/pkg/store/redis"
	"github.com/redis/go-redis/v9"
	// "github.com/jin06/binlogo/v2/pkg/store/store_redis"
)

const (
	registerDuration = time.Second * 5
)

// NodePrefix returns etcd prefix of node
func NodePrefix() string {
	return etcdclient.Prefix() + "/node/info"
}

// NodeRegisterPrefix returns etcd prefix of register node
func NodeRegisterPrefix() string {
	return etcdclient.Prefix() + "/cluster/register"
}

func RegisterNode(ctx context.Context, n *node.Node) (bool, error) {
	return store_redis.GetClient().SetNX(ctx, n.RegisterName(), n.Name, registerDuration).Result()
}

func LeaseNode(ctx context.Context, n *node.Node) (bool, error) {
	err := store_redis.GetClient().SetEx(ctx, n.RegisterName(), n.Name, registerDuration).Err()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
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
func AllNodesMap() (mapping map[string]*node.Node, err error) {
	list, err := AllNodes(context.Background())
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
func AllRegisterNodesMap() (mapping map[string]*node.Node, err error) {
	list, err := AllNodes(context.Background())
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
	res, err = AllNodesMap()
	if err != nil {
		return
	}
	regMap, err := AllRegisterNodesMap()
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
