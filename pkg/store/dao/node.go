package dao

import (
	"context"
	"encoding/json"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"go.etcd.io/etcd/clientv3"
)

func CreateNode(n *node.Node) (err error) {
	key := etcd.Prefix() + "/nodes" + n.Name
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
	key := etcd.Prefix() + "/nodes/" + n.Name
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

func GetNode(name string) (n *node.Node, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), etcd.E.Timeout)
	defer cancel()
	key := etcd.Prefix() + "/nodes" + name
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
	key := etcd.E.Prefix + "/nodes"
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
	for k,v := range list {
		if v.Status	== node.STATUS_OFF {
			list = append(list[:k], list[k:]...)
		}
	}
	return
}
