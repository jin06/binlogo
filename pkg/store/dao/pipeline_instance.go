package dao

import (
	"context"
	"time"

	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
)

// GetInstance get a pipeline instance from etcd
// get by pipeline name pipeName
func GetInstance(ctx context.Context, pname string) (ins *pipeline.Instance, err error) {
	return myDao.GetInstance(ctx, pname)
}

func RegisterInstance(ctx context.Context, ins *pipeline.Instance, exp time.Duration) error {
	return myDao.RegisterInstance(ctx, ins, exp)
}

func LeaseInstance(ctx context.Context, pipe string, exp time.Duration) error {
	return myDao.LeaseInstance(ctx, pipe, exp)
}

// AllInstance returns all pipeline instances in array form
func AllInstance(ctx context.Context) (all []*pipeline.Instance, err error) {
	return myDao.AllInstance(ctx)
}

// AllInstanceMap returns all pipeline instances in map form
func AllInstanceMap(ctx context.Context) (all map[string]*pipeline.Instance, err error) {
	return myDao.AllInstanceMap(ctx)
}

func UnRegisterInstance(ctx context.Context, pipe string, n string) error {
	return myDao.UnRegisterInstance(ctx, pipe, n)
}
