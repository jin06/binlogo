package dao

import (
	"context"

	"github.com/jin06/binlogo/v2/pkg/store/model"
)

// Update create or update event
func UpdateEvent(ctx context.Context, e *model.Event) error {
	return myDao.UpdateEvent(ctx, e)
}
