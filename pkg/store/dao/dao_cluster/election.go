package dao_cluster

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/jin06/binlogo/pkg/etcdclient"
)

// ElectionPrefix returns etcd prefix
func ElectionPrefix() string {
	return etcdclient.Prefix() + "/cluster/election"
}

// LeaderNode return current name of current leader node
func LeaderNode() (nodeName string, err error) {
	key := ElectionPrefix()
	resp, err := etcdclient.Default().Get(context.Background(), key, clientv3.WithFirstCreate()...)
	if err != nil {
		return
	}
	if len(resp.Kvs) == 0 {
		return
	}
	nodeName = string(resp.Kvs[0].Value)
	return
}

// AllElections returns all nodes in the election
func AllElections() (res []map[string]interface{}, err error) {
	key := ElectionPrefix()
	resp, err := etcdclient.Default().Get(context.Background(), key, clientv3.WithPrefix())
	if err != nil {
		return
	}
	res = []map[string]interface{}{}

	for _, v := range resp.Kvs {
		item := map[string]interface{}{}
		item["create_revision"] = v.CreateRevision
		item["node"] = string(v.Value)
		res = append(res, item)
	}
	return
}
