package dao_node

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/coreos/etcd/clientv3"
	"github.com/jin06/binlogo/pkg/etcdclient"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/sirupsen/logrus"
)

// CapacityPrefix returns prefix key of etcd for capacity
func CapacityPrefix() string {
	return etcdclient.Prefix() + "/node/capacity"
}

// UpdateCapacity update capacity data to etcd
func UpdateCapacity(cap *node.Capacity) (err error) {
	if cap.NodeName == "" {
		err = errors.New("empty name")
		return
	}
	key := CapacityPrefix() + "/" + cap.NodeName
	b, err := json.Marshal(cap)
	if err != nil {
		return
	}
	_, err = etcdclient.Default().Put(context.TODO(), key, string(b))
	return
}

// CapacityMap returns capacity data in map form
func CapacityMap() (mapping map[string]*node.Capacity, err error) {
	key := CapacityPrefix()
	res, err := etcdclient.Default().Get(context.TODO(), key, clientv3.WithPrefix())
	if err != nil {
		return
	}
	if len(res.Kvs) == 0 {
		return
	}
	mapping = map[string]*node.Capacity{}
	for _, v := range res.Kvs {
		ele := &node.Capacity{}
		er := json.Unmarshal(v.Value, ele)
		if er != nil {
			logrus.Error(er)
			continue
		}
		var nodeName string
		_, er = fmt.Sscanf(string(v.Key), key+"/%s", &nodeName)
		if er != nil {
			logrus.Error(er)
			continue
		}
		if nodeName != "" {
			mapping[nodeName] = ele
		}
	}

	return
}

// DeleteCapacity delete capacity in etcd
func DeleteCapacity(nodeName string) (ok bool, err error) {
	if nodeName == "" {
		err = errors.New("empty name")
		return
	}
	key := CapacityPrefix() + "/" + nodeName
	res, err := etcdclient.Default().Delete(context.Background(), key)
	if err != nil {
		return
	}
	if res.Deleted > 0 {
		ok = true
	}
	return
}
