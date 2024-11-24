package dao

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jin06/binlogo/v2/pkg/etcdclient"
	"github.com/jin06/binlogo/v2/pkg/store/model"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
	store_redis "github.com/jin06/binlogo/v2/pkg/store/redis"
	"github.com/redis/go-redis/v9"

	// "github.com/jin06/binlogo/v2/pkg/store/store_redis"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
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

// CreateNode create a node in etcd
func CreateNode(n *node.Node) (err error) {
	key := NodePrefix() + "/" + n.Name
	v, err := json.Marshal(n)
	if err != nil {
		return
	}
	_, err = etcdclient.Default().Put(context.Background(), key, string(v))
	return
}

// CreateNodeIfNotExist create a node in etcd if that not exist
func CreateNodeIfNotExist(ctx context.Context, n *node.Node) error {
	if n == nil {
		return errors.New("empty node")
	}
	if len(n.Name) == 0 {
		return errors.New("empty node name")
	}
	_, err := store_redis.Default.Create(ctx, n)
	return err
}

// UpdateNode update a node data in etcd
func UpdateNode(ctx context.Context, name string, opts ...node.NodeOption) (bool, error) {
	if len(name) == 0 {
		return false, errors.New("empty node name")
	}
	values := model.Values{}
	for _, v := range opts {
		v(values)
	}
	return store_redis.Default.UpdateField(ctx, &node.Node{Name: name}, values)
}

// GetNode get node data from etcd
func GetNode(ctx context.Context, name string) (*node.Node, error) {
	n := &node.Node{}
	has, err := store_redis.Default.Get(ctx, n)
	if err != nil {
		return nil, err
	}
	if has {
		return n, nil
	}
	return nil, nil
}

// AllNodes return all node data from etcd
func AllNodes() (list []*node.Node, err error) {
	return _allNodes(NodePrefix() + "/")
}

func _allNodes(key string) (list []*node.Node, err error) {
	list = []*node.Node{}
	//key := NodePrefix() + "/"
	res, err := etcdclient.Default().Get(context.TODO(), key, clientv3.WithPrefix())
	if err != nil {
		return
	}
	if len(res.Kvs) == 0 {
		return
	}
	for _, v := range res.Kvs {
		ele := &node.Node{}
		er := ele.Unmarshal(v.Value)
		if er != nil {
			logrus.Error(er)
			continue
		}
		list = append(list, ele)
	}
	return
}

// ALLRegisterNodes returns all register nodes
func ALLRegisterNodes() (list []*node.Node, err error) {
	key := NodeRegisterPrefix()
	return _allNodes(key)
}

// AllNodesMap  returns all register nodes in map form
func AllNodesMap() (mapping map[string]*node.Node, err error) {
	list, err := AllNodes()
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
	list, err := ALLRegisterNodes()
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

// DeleteNode delete a node by node name
func DeleteNode(name string) (ok bool, err error) {
	if name == "" {
		err = errors.New("empty node name")
		return
	}
	key := NodePrefix() + "/" + name
	res, err := etcdclient.Default().Delete(context.Background(), key)
	if err != nil {
		return
	}
	if res.Deleted > 0 {
		ok = true
	}
	return
}
