package dao_pipe

import (
	"context"
	"encoding/json"

	"github.com/jin06/binlogo/v2/pkg/etcdclient"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// InstancePrefix returns etcd prefix of pipeline instance
func InstancePrefix() string {
	return etcdclient.Prefix() + "/pipeline/instance"
}

// GetInstance get a pipeline instance from etcd
// get by pipeline name pipeName
func GetInstance(pipeName string) (ins *pipeline.Instance, err error) {
	key := InstancePrefix() + "/" + pipeName
	res, err := etcdclient.Default().Get(context.Background(), key)
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

// AllInstance returns all pipeline instances in array form
func AllInstance() (all []*pipeline.Instance, err error) {
	all = []*pipeline.Instance{}
	key := InstancePrefix() + "/"
	res, err := etcdclient.Default().Get(context.Background(), key, clientv3.WithPrefix())
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
		all = append(all, item)
	}
	return
}

// AllInstanceMap returns all pipeline instances in map form
func AllInstanceMap() (all map[string]*pipeline.Instance, err error) {
	all = map[string]*pipeline.Instance{}
	key := InstancePrefix() + "/"
	res, err := etcdclient.Default().Get(context.Background(), key, clientv3.WithPrefix())
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
