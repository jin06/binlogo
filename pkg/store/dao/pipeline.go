package dao

import (
	"context"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"go.etcd.io/etcd/clientv3"
)

func CreatePipeline(d *pipeline.Pipeline, opts ...clientv3.OpOption) (ok bool, err error) {
	return etcd.Create(d, opts...)
}

func UpdatePipeline(d *pipeline.Pipeline) (bool, error) {
	return etcd.Update(d)
}

func AllPipelines() (list []*pipeline.Pipeline, err error) {
	list = []*pipeline.Pipeline{}
	key := etcd.Prefix() + "/pipeline"
	res, err := etcd.E.Client.Get(context.Background(), key, clientv3.WithPrefix(), clientv3.WithFromKey())
	if err != nil {
		return
	}
	if len(res.Kvs) == 0 {
		return
	}
	for _, v := range res.Kvs {
		ele := &pipeline.Pipeline{
		}
		er := ele.Unmarshal(v.Value)
		if er != nil {
			blog.Error(er)
			continue
		}
		list = append(list, ele)
	}
	return
}
