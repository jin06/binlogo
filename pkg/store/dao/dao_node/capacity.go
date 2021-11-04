package dao_node

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/jin06/binlogo/pkg/store/model/node"
)

func CapacityPrefix() string {
	return etcd.Prefix() + "/node/capacity"
}

func UpdateCapacity(cap *node.Capacity, args ...Option) (err error) {
	opts := GetOptions(args...)
	key := opts.Key
	if key == "" {
		if opts.NodeName != "" {
			key = CapacityPrefix() + "/" + opts.NodeName
		}
	}
	if key == "" {
		err = errors.New("empty key")
		return
	}
	ctx, cancel := context.WithTimeout(context.TODO(), etcd.E.Timeout)
	defer cancel()
	b, err := json.Marshal(cap)
	if err != nil {
		return
	}
	_, err = etcd.E.Client.Put(ctx, key, string(b))
	return
}

func CapacityMap() (mapping map[string]*node.Capacity, err error) {
	key := CapacityPrefix()
	res, err := etcd.E.Client.Get(context.TODO(), key, clientv3.WithPrefix())
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
			blog.Error(er)
			continue
		}
		var nodeName string
		_, er = fmt.Sscanf(string(v.Key), key+"/%s", &nodeName)
		if er != nil {
			blog.Error(er)
			continue
		}
		if nodeName != "" {
			mapping[nodeName] = ele
		}
	}

	return
}
