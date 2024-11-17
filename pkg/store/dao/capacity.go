package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jin06/binlogo/v2/pkg/etcdclient"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
	store_redis "github.com/jin06/binlogo/v2/pkg/store/redis"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// CapacityPrefix returns prefix key of etcd for capacity
func CapacityPrefix() string {
	return etcdclient.Prefix() + "/node/capacity"
}

// UpdateCapacity update capacity data to etcd
func UpdateCapacity(ctx context.Context, cap *node.Capacity) (bool, error) {
	return store_redis.Default.Update(ctx, cap)
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
func DeleteCapacity(ctx context.Context, nodeName string) (bool, error) {
	return store_redis.Default.Delete(ctx, &node.Capacity{NodeName: nodeName})
}
