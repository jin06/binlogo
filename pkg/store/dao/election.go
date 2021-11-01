package dao

import (
	"context"
	"github.com/jin06/binlogo/pkg/store/etcd"
	"go.etcd.io/etcd/clientv3"
)

func LeaderNode() (nodeName string, err error) {
	key := etcd.Prefix() + "/election"
	resp, err := etcd.E.Client.Get(context.TODO(), key, clientv3.WithFirstCreate()...)
	if err != nil {
		return
	}
	if len(resp.Kvs) == 0 {
		return
	}
	nodeName = string(resp.Kvs[0].Value)
	return
}
