package dao

import (
	"context"

	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
)

// UpdatePosition update pipeline record position in etcd
func UpdateRecord(ctx context.Context, p *pipeline.RecordPosition) error {
	return myDao.UpdateRecord(ctx, p)
}

// UpdateRecordSafe  update pipeline record position in etcd in safe mode.
// there will be version judgment when updating
func UpdateRecordSafe(ctx context.Context, pipeName string, opts ...pipeline.OptionRecord) error {
	return myDao.UpdateRecordSafe(ctx, pipeName, opts...)
}

// GetRecord get pipeline record position from etcd
func GetRecord(ctx context.Context, pipeName string) (*pipeline.RecordPosition, error) {
	return myDao.GetRecord(ctx, pipeName)
}
