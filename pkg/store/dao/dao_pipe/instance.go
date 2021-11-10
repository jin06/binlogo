package dao_pipe

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/jin06/binlogo/pkg/etcd_client"
	"github.com/sirupsen/logrus"
)

func InstancePrefix() string {
	return etcd_client.Prefix() + "/pipeline/instance"
}

func GetInstance(pipeName string) (nodeName string, err error) {
	key := InstancePrefix() + "/" + pipeName
	res, err := etcd_client.Default().Get(context.Background(), key)
	if err != nil {
		return
	}
	if len(res.Kvs) == 0 {
		return
	}
	err = json.Unmarshal(res.Kvs[0].Value, &nodeName)
	if err != nil {
		return
	}
	return
}

func AllInstance() (all map[string]string, err error) {
	all = map[string]string{}
	key := InstancePrefix()
	res, err := etcd_client.Default().Get(context.Background(), key, clientv3.WithPrefix())
	if err != nil {
		return
	}
	if len(res.Kvs) == 0 {
		return
	}
	for _, v := range res.Kvs {
		var item string
		er := json.Unmarshal(v.Value, &item)
		if er != nil {
			logrus.Error(er)
			continue
		}
		var index string
		_, er = fmt.Sscanf(string(v.Key), InstancePrefix()+"/%s", &index)
		if er != nil || index == "" {
			continue
		}
		all[index] = item
	}
	return
}
