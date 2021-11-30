package dao_node

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/coreos/etcd/clientv3"
	"github.com/jin06/binlogo/pkg/etcdclient"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/sirupsen/logrus"
)

// NodePrefix returns etcd prefix of node
func NodePrefix() string {
	return etcdclient.Prefix() + "/node/info"
}

// NodeRegisterPrefix returns etcd prefix of register node
func NodeRegisterPrefix() string {
	return etcdclient.Prefix() + "/cluster/register"
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
func CreateNodeIfNotExist(n *node.Node) (err error) {
	key := NodePrefix() + "/" + n.Name
	b, err := json.Marshal(n)
	if err != nil {
		return
	}
	ctx := context.Background()
	txn := etcdclient.Default().Txn(ctx).If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0))
	txn = txn.Then(clientv3.OpPut(key, string(b)))
	resp, err := txn.Commit()

	if err != nil {
		logrus.Error(err)
		logrus.Error(resp.Succeeded)
	}
	return
}

// UpdateNode update a node data in etcd
func UpdateNode(nodeName string, opts ...node.NodeOption) (ok bool, err error) {
	if nodeName == "" {
		err = errors.New("empty node name")
		return
	}
	key := NodePrefix() + "/" + nodeName
	res, err := etcdclient.Default().Get(context.TODO(), key)
	if err != nil {
		return
	}
	if len(res.Kvs) <= 0 {
		err = errors.New("empty node")
		return
	}
	revision := int64(0)
	revision = res.Kvs[0].CreateRevision
	n := &node.Node{}
	err = json.Unmarshal(res.Kvs[0].Value, n)
	if err != nil {
		return
	}
	for _, v := range opts {
		v(n)
	}

	txn := etcdclient.Default().Txn(context.TODO()).If(clientv3.Compare(clientv3.CreateRevision(key), "=", revision))
	txn = txn.Then(clientv3.OpPut(key, n.Val()))
	resp, err := txn.Commit()
	if err != nil {
		return
	}
	ok = resp.Succeeded
	return
}

// GetNode get node data from etcd
func GetNode(name string) (n *node.Node, err error) {
	key := NodePrefix() + "/" + name
	res, err := etcdclient.Default().Get(context.Background(), key)
	if err != nil {
		return
	}
	if len(res.Kvs) == 0 {
		return
	}
	n = &node.Node{}
	err = n.Unmarshal(res.Kvs[0].Value)
	return
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
		err = errors.New("empty name")
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
