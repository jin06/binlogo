package dao

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jin06/binlogo/v2/pkg/etcdclient"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
	store_redis "github.com/jin06/binlogo/v2/pkg/store/redis"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// StatusPrefix returns etcd prefix of node status
func StatusPrefix() string {
	return etcdclient.Prefix() + "/node/status"
}

// CreateOrUpdateStatus crate of update status in etcd
// create if not exist
func CreateOrUpdateStatus(nodeName string, opts ...node.StatusOption) (ok bool, err error) {
	if nodeName == "" {
		err = errors.New("empty node name")
		return
	}
	key := StatusPrefix() + "/" + nodeName
	res, err := etcdclient.Default().Get(context.TODO(), key)
	if err != nil {
		return
	}
	revision := int64(0)
	s := &node.Status{}
	s.NodeName = nodeName
	if len(res.Kvs) > 0 {
		revision = res.Kvs[0].CreateRevision
		err = json.Unmarshal(res.Kvs[0].Value, &s)
		if err != nil {
			return
		}
	}
	for _, v := range opts {
		v(s)
	}

	txn := etcdclient.Default().Txn(context.TODO()).If(clientv3.Compare(clientv3.CreateRevision(key), "=", revision))
	b, _ := json.Marshal(s)
	txn = txn.Then(clientv3.OpPut(key, string(b)))
	resp, err := txn.Commit()
	if err != nil {
		return
	}
	ok = resp.Succeeded
	return
}

// GetStatus get node status from etcd
func GetStatus(ctx context.Context, nodeName string) (s *node.Status, err error) {
	store_redis.Default.Get(ctx, &node.Status{NodeName: nodeName})
	key := StatusPrefix() + "/" + nodeName
	res, err := etcdclient.Default().Get(context.TODO(), key)
	if err != nil {
		return
	}
	if len(res.Kvs) == 0 {
		return
	}
	s = &node.Status{}
	err = json.Unmarshal(res.Kvs[0].Value, s)
	return
}

// CreateStatusIfNotExist create status if not exist
func CreateStatusIfNotExist(ctx context.Context, nodeName string, opts ...node.StatusOption) (ok bool, err error) {
	return myDao.CreateOrUpdateStatus(ctx, nodeName, opts...)
}

// StatusMap returns all node status in map form
func StatusMap(ctx context.Context) (mapping map[string]*node.Status, err error) {
	return myDao.StatusMap(ctx)
}

// DeleteStatus delete node status in etcd
func DeleteStatus(name string) (ok bool, err error) {
	if name == "" {
		err = errors.New("empty name")
		return
	}
	key := StatusPrefix() + "/" + name
	res, err := etcdclient.Default().Delete(context.Background(), key)
	if err != nil {
		return
	}
	if res.Deleted > 0 {
		ok = true
	}
	return
}
