package dao

import (
	"context"

	"github.com/jin06/binlogo/v2/pkg/store/model"
)

// GetPipelineBind get pipeline bind from etcd
func GetPipelineBind(ctx context.Context) (*model.PipelineBind, error) {
	return myDao.GetPipelineBind(ctx)
}

// UpdatePipelineBindIfNotExist update pipeline bind if the pipeline not exist in pipeline bind
func UpdatePipelineBindIfNotExist(ctx context.Context, pName string, nName string) (err error) {
	return myDao.UpdatePipelineBindIfNotExist(ctx, pName, nName)
}

// UpdatePipelineBind update pipeline bind ignore exists
func UpdatePipelineBind(ctx context.Context, pName string, nName string) (ok bool, err error) {
	return myDao.UpdatePipelineBind(ctx, pName, nName)
}

// DeletePipelineBind delete pipeline bind in etcd
func DeletePipelineBind(ctx context.Context, pName string) (ok bool, err error) {
	return myDao.DeletePipelineBind(ctx, pName)
}
