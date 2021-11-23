package dao_node

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/coreos/etcd/clientv3"
	"github.com/jin06/binlogo/pkg/etcd_client"
	"github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/sirupsen/logrus"
)

// StatusPrefix returns etcd prefix of node status
func StatusPrefix() string {
	return etcd_client.Prefix() + "/node/status"
}

// CreateOrUpdateStatus crate of update status in etcd
// create if not exist
func CreateOrUpdateStatus(nodeName string, opts ...node.StatusOption) (ok bool, err error) {
	key := StatusPrefix() + "/" + nodeName
	res, err := etcd.E.Client.Get(context.TODO(), key)
	if err != nil {
		return
	}
	revision := int64(0)
	s := &node.Status{}
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

	txn := etcd.E.Client.Txn(context.TODO()).If(clientv3.Compare(clientv3.CreateRevision(key), "=", revision))
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
func GetStatus(nodeName string) (s *node.Status, err error) {
	key := NodePrefix() + "/" + nodeName
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

// CreateStatusIfNotExist create status if not exist
func CreateStatusIfNotExist(n *node.Status) (err error) {
	if n.NodeName == "" {
		return errors.New("empty name")
	}
	key := StatusPrefix() + "/" + n.NodeName
	b, err := json.Marshal(n)
	if err != nil {
		return
	}
	txn := etcd.E.Client.Txn(context.TODO()).If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0))
	txn = txn.Then(clientv3.OpPut(key, string(b)))
	resp, err := txn.Commit()

	if err != nil {
		logrus.Error(err)
		logrus.Error(resp.Succeeded)
	}
	return
}

// StatusMap returns all node status in map form
func StatusMap() (mapping map[string]*node.Status, err error) {
	key := StatusPrefix()
	res, err := etcd.E.Client.Get(context.TODO(), key, clientv3.WithPrefix())
	if err != nil {
		return
	}
	if len(res.Kvs) == 0 {
		return
	}
	mapping = map[string]*node.Status{}
	for _, v := range res.Kvs {
		ele := &node.Status{}
		er := json.Unmarshal(v.Value, ele)
		if er != nil {
			logrus.Error(er)
			continue
		}
		mapping[ele.NodeName] = ele
	}

	return
}

// DeleteStatus delete node status in etcd
func DeleteStatus(name string) (err error) {
	if name == "" {
		return errors.New("empty name")
	}
	key := StatusPrefix() + "/" + name
	_, err = etcd_client.Default().Delete(context.Background(), key)
	if err != nil {
		return
	}
	return
}
