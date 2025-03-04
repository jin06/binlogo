package pipeline

import (
	"context"
	"errors"

	"github.com/jin06/binlogo/v2/pkg/store/dao"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
)

// GetItemByName get item by pipeline name
func GetItemByName(ctx context.Context, name string) (item *Item, err error) {
	if name == "" {
		err = errors.New("empty name")
		return
	}
	pipe, err := dao.GetPipeline(ctx, name)
	if err != nil {
		return
	}
	if pipe == nil {
		return
	}
	item = &Item{
		Pipeline: pipe,
	}
	err = CompleteInfo(item)
	if err != nil {
		return
	}
	return
}

// PipeStatus get pipeline status by pipeline name
func PipeStatus(ctx context.Context, name string) (status pipeline.Status, err error) {
	pipe, err := dao.GetPipeline(ctx, name)
	if err != nil {
		return
	}
	if pipe == nil {
		return
	}
	status = pipe.Status
	return
}
