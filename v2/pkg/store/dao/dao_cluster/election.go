package dao_cluster

import (
	"context"
	"errors"

	"github.com/jin06/binlogo/v2/pkg/etcdclient"
	clientv3 "go.etcd.io/etcd/client/v3"
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

// GetElection get election value
func GetElection(leaseId string) (res *string, err error) {
	if leaseId == "" {
		err = errors.New("empty lease id")
		return
	}
	key := ElectionPrefix() + "/" + leaseId
	resp, err := etcdclient.Default().Get(context.Background(), key)
	if err != nil {
		return
	}
	if len(resp.Kvs) == 0 {
		return
	}
	val := string(resp.Kvs[0].Value)
	res = &val
	return
}
