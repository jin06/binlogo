package dao_pipe

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"github.com/jin06/binlogo/pkg/etcd_client"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/sirupsen/logrus"
)

func InstancePrefix() string {
	return etcd_client.Prefix() + "/pipeline/instance"
}

func GetInstance(pipeName string) (ins *pipeline.Instance, err error) {
	key := InstancePrefix() + "/" + pipeName
	res, err := etcd_client.Default().Get(context.Background(), key)
	if err != nil {
		return
	}
	if len(res.Kvs) == 0 {
		return
	}
	ins = &pipeline.Instance{}
	err = json.Unmarshal(res.Kvs[0].Value, ins)
	if err != nil {
		return
	}
	return
}

func AllInstance() (all map[string]*pipeline.Instance, err error) {
	all = map[string]*pipeline.Instance{}
	key := InstancePrefix() + "/"
	res, err := etcd_client.Default().Get(context.Background(), key, clientv3.WithPrefix())
	if err != nil {
		return
	}
	if len(res.Kvs) == 0 {
		return
	}
	for _, v := range res.Kvs {
		item := &pipeline.Instance{}
		er := json.Unmarshal(v.Value, item)
		if er != nil {
			logrus.Error(er)
			continue
		}
		all[item.PipelineName] = item
	}
	return
}

