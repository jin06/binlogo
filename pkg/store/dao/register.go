package dao

import (
	"context"
	"encoding/json"
	"github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/jin06/binlogo/pkg/store/model/register"
	"go.etcd.io/etcd/clientv3"
)

func RegRNode(n *register.RegisterNode,leaseID clientv3.LeaseID)(err error) {
	key := etcd.Prefix() + "/register/" + n.Name
	b, err := json.Marshal(n)
	if err != nil {
		return
	}
	_, err =  etcd.E.Client.Put(context.Background(), key, string(b), clientv3.WithLease(leaseID))
	if err != nil {
		return
	}
	return
}

func GetRNode(name string) (n *register.RegisterNode, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), etcd.E.Timeout)
	defer cancel()
	key := etcd.Prefix() + "/register/" + name
	res, err := etcd.E.Client.Get(ctx, key)
	if err != nil {
		return
	}
	if len(res.Kvs) == 0 {
		return
	}
	n = &register.RegisterNode{}
	err = json.Unmarshal(res.Kvs[0].Value, n)
	return
}
