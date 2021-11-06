package dao_cluster

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/jin06/binlogo/pkg/store/model/register"
)

func RegisterPrefix() string {
	return etcd.Prefix() + "/cluster/register"
}

func RegRNode(n *register.RegisterNode, leaseID clientv3.LeaseID) (err error) {
	key := RegisterPrefix() + "/" + n.Name
	b, err := json.Marshal(n)
	if err != nil {
		return
	}
	_, err = etcd.E.Client.Put(context.Background(), key, string(b), clientv3.WithLease(leaseID))
	if err != nil {
		return
	}
	return
}

func GetRNode(name string) (n *register.RegisterNode, err error) {
	key := RegisterPrefix() + "/" + name
	res, err := etcd.E.Client.Get(context.Background(), key)
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
