package dao

import (
	"context"

	"github.com/jin06/binlogo/v2/pkg/etcdclient"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
)

// InstancePrefix returns etcd prefix of pipeline instance
func InstancePrefix() string {
	return etcdclient.Prefix() + "/pipeline/instance"
}

// GetInstance get a pipeline instance from etcd
// get by pipeline name pipeName
func GetInstance(ctx context.Context, pname string) (ins *pipeline.Instance, err error) {
	return myDao.GetInstance(ctx, pname)
}

// AllInstance returns all pipeline instances in array form
func AllInstance(ctx context.Context) (all []*pipeline.Instance, err error) {
	return myDao.AllInstance(ctx)
}

// AllInstanceMap returns all pipeline instances in map form
func AllInstanceMap(ctx context.Context) (all map[string]*pipeline.Instance, err error) {
	return myDao.AllInstanceMap(ctx)
}
