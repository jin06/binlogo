package dao_node

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/coreos/etcd/clientv3"
	"github.com/jin06/binlogo/pkg/etcd_client"
	"github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/sirupsen/logrus"
)

func NodePrefix() string {
	return etcd.Prefix() + "/node/info"
}

func NodeRegisterPrefix() string {
	return etcd_client.Prefix() + "/cluster/register"
}

func CreateNode(n *node.Node) (err error) {
	key := NodePrefix() + "/" + n.Name
	ctx, cancel := context.WithTimeout(context.TODO(), etcd.E.Timeout)
	defer cancel()
	v, err := json.Marshal(n)
	if err != nil {
		return
	}
	_, err = etcd.E.Client.Put(ctx, key, string(v))
	return
}

func CreateNodeIfNotExist(n *node.Node) (err error) {
	key := NodePrefix() + "/" + n.Name
	ctx, cancel := context.WithTimeout(context.TODO(), etcd.E.Timeout)
	defer cancel()
	b, err := json.Marshal(n)
	if err != nil {
		return
	}
	txn := etcd.E.Client.Txn(ctx).If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0))
	txn = txn.Then(clientv3.OpPut(key, string(b)))
	resp, err := txn.Commit()

	if err != nil {
		logrus.Error(err)
		logrus.Error(resp.Succeeded)
	}
	return
}

func UpdateNode(nodeName string, opts ...node.NodeOption) (ok bool, err error) {
	if nodeName == "" {
		err = errors.New("empty node name")
		return
	}
	key := NodePrefix() + "/" + nodeName
	res, err := etcd_client.Default().Get(context.TODO(), key)
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

	txn := etcd_client.Default().Txn(context.TODO()).If(clientv3.Compare(clientv3.CreateRevision(key), "=", revision))
	txn = txn.Then(clientv3.OpPut(key, n.Val()))
	resp, err := txn.Commit()
	if err != nil {
		return
	}
	ok = resp.Succeeded
	return
}

func GetNode(name string) (n *node.Node, err error) {
	ctx, cancel := context.WithTimeout(context.TODO(), etcd.E.Timeout)
	defer cancel()
	key := NodePrefix() + "/" + name
	res, err := etcd.E.Client.Get(ctx, key)
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

func AllNodes() (list []*node.Node, err error) {
	return _allNodes(NodePrefix() + "/")
}
func _allNodes(key string) (list []*node.Node, err error) {
	list = []*node.Node{}
	//key := NodePrefix() + "/"
	res, err := etcd.E.Client.Get(context.TODO(), key, clientv3.WithPrefix())
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

func ALLRegisterNodes() (list []*node.Node, err error) {
	key := NodeRegisterPrefix()
	return _allNodes(key)
}

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

func AllWorkNodesMap() (res map[string]*node.Node, err error) {
	res, err = AllNodesMap()
	if err != nil {
		return
	}
	regMap, err :=  AllRegisterNodesMap()
	if err != nil {
		return
	}
	for k ,_ := range res {
		if _, ok := regMap[k]; !ok {
			delete(res, k)
		}
	}
	return
}
