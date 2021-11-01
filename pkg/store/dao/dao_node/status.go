package dao_node

import (
	"context"
	"encoding/json"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"go.etcd.io/etcd/clientv3"
)

func nodePrefix() string {
	return etcd.Prefix() + "/nodes"
}
func statsPrefix() string {
	return etcd.Prefix() + "/status"
}

func CreateOrUpdateStatus(nodeName string, opts ...Option) (ok bool, err error) {
	options := GetOptions(opts...)
	key := statsPrefix() + "/" + nodeName
	res , err := etcd.E.Client.Get(context.TODO(), key)
	if err != nil {
		return
	}
	revision := int64(0)
	s := node.Status{}
	if len(res.Kvs) > 0 {
		revision = res.Kvs[0].CreateRevision
		err = json.Unmarshal(res.Kvs[0].Value, &s )
		if err != nil {
			return
		}
	}
	if options.StatusReady != nil {
		s.Ready = options.StatusReady.val
	}
	txn := etcd.E.Client.Txn(context.TODO()).If(clientv3.Compare(clientv3.CreateRevision(key), "=", revision))
	b, _  := json.Marshal(&s)
	txn = txn.Then(clientv3.OpPut(key, string(b)))
	resp , err := txn.Commit()
	if err != nil {
		return
	}
	ok = resp.Succeeded
	return
}

func GetStatus(nodeName string) (s *node.Status, err error) {
	key := nodePrefix() + "/" + nodeName
	res, err := etcd.E.Client.Get(context.TODO(), key)
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

func CreateStatusIfNotExist(nodeName string, n *node.Status) (err error) {
	key := etcd.Prefix() + "/status/" + nodeName
	b, err := json.Marshal(n)
	if err != nil {
		return
	}
	txn := etcd.E.Client.Txn(context.TODO()).If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0))
	txn = txn.Then(clientv3.OpPut(key, string(b)))
	resp, err := txn.Commit()

	if err != nil {
		blog.Error(err)
		blog.Error(resp.Succeeded)
	}
	return
}
