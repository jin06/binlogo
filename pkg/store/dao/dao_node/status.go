package dao_node

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/sirupsen/logrus"
)


func StatsPrefix() string {
	return etcd.Prefix() + "/node/status"
}

func CreateOrUpdateStatus(nodeName string, opts ...Option) (ok bool, err error) {
	options := GetOptions(opts...)
	key := StatsPrefix() + "/" + nodeName
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

func CreateStatusIfNotExist(nodeName string, n *node.Status) (err error) {
	key := StatsPrefix() + "/" + nodeName
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

func StatusMap() (mapping map[string]*node.Status, err error) {
	key := StatsPrefix()
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
