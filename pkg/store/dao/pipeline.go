package dao

import (
	"github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"go.etcd.io/etcd/clientv3"
)

func CreatePipeline(d *pipeline.Pipeline, opts ...clientv3.OpOption) (ok bool, err error) {
	return etcd.Create(d, opts...)
}

func UpdatePipeline(d * pipeline.Pipeline) (bool, error) {
	return etcd.Update(d)
}

