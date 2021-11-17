package dao_cluster

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/jin06/binlogo/pkg/etcd_client"
)

func ElectionPrefix() string {
	return etcd_client.Prefix() + "/cluster/election"
}
func LeaderNode() (nodeName string, err error) {
	key := ElectionPrefix()
	resp, err := etcd_client.Default().Get(context.Background(), key, clientv3.WithFirstCreate()...)
	if err != nil {
		return
	}
	if len(resp.Kvs) == 0 {
		return
	}
	nodeName = string(resp.Kvs[0].Value)
	return
}

func AllElections() (res []map[string]interface{}, err error) {
	key := ElectionPrefix()
	resp, err := etcd_client.Default().Get(context.Background(), key, clientv3.WithPrefix())
	if err != nil {
		return
	}
	res = []map[string]interface{}{}

	for _,v := range resp.Kvs {
		item := map[string]interface{}{}
		item["create_revision"] = v.CreateRevision
		item["node"] = string(v.Value)
		fmt.Println(v.ProtoMessage)
		res = append(res, item)
	}
	return
}
