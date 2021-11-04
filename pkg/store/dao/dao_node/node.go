package dao_node

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/coreos/etcd/clientv3"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/jin06/binlogo/pkg/store/model/node"
)

func NodePrefix() string {
	return etcd.Prefix() + "/node/info"
}

func CreateNode(n *node.Node) (err error) {
	key := NodePrefix() + "/" + n.Name
	ctx, cancel := context.WithTimeout(context.Background(), etcd.E.Timeout)
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
	ctx, cancel := context.WithTimeout(context.Background(), etcd.E.Timeout)
	defer cancel()
	b, err := json.Marshal(n)
	if err != nil {
		return
	}
	txn := etcd.E.Client.Txn(ctx).If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0))
	txn = txn.Then(clientv3.OpPut(key, string(b)))
	resp, err := txn.Commit()

	if err != nil {
		blog.Error(err)
		blog.Error(resp.Succeeded)
	}
	return
}

func UpdateNode(nodeName string, opts ...Option) (ok bool, err error) {
	if nodeName == "" {
		err = errors.New("empty node name")
		return
	}
	key := NodePrefix() + "/" + nodeName
	res, err := etcd.E.Client.Get(context.TODO(), key)
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
	options := GetOptions(opts...)
	if options.NodeIP != nil {
		n.IP = options.NodeIP
	}
	if options.NodeVersion != "" {
		n.Version = options.NodeVersion
	}
	txn := etcd.E.Client.Txn(context.TODO()).If(clientv3.Compare(clientv3.CreateRevision(key), "=", revision))
	txn = txn.Then(clientv3.OpPut(key, n.Val()))
	resp, err := txn.Commit()
	if err != nil {
		return
	}
	ok = resp.Succeeded
	return
}

func GetNode(name string) (n *node.Node, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), etcd.E.Timeout)
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
	list = []*node.Node{}
	key := NodePrefix() + "/"
	res, err := etcd.E.Client.Get(context.Background(), key, clientv3.WithPrefix())
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
			blog.Error(er)
			continue
		}
		list = append(list, ele)
	}
	return
}

func AllWorkNodes() (list []*node.Node, err error) {
	list, err = AllNodes()
	if err != nil {
		return
	}
	//for k, v := range list {
	//	if v.Status == node.STATUS_OFF {
	//		list = append(list[:k], list[k:]...)
	//	}
	//}
	return
}
