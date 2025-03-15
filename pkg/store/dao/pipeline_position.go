package dao

import (
	"context"

	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
)

// UpdatePosition update pipeline position in etcd
func UpdatePosition(ctx context.Context, p *pipeline.Position) (err error) {
	return myDao.UpdatePosition(ctx, p)
}

// GetPosition get pipeline position from etcd
func GetPosition(ctx context.Context, pipe string) (p *pipeline.Position, err error) {
	return myDao.GetPosition(ctx, pipe)
}

// DeletePosition delete pipeline positon by pipeline name in etcd
func DeletePosition(ctx context.Context, name string) error {
	return myDao.DeletePosition(ctx, name)
}
