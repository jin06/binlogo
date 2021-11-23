package pipeline

import (
	"errors"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
)

// GetItemByName get item by pipeline name
func GetItemByName(name string) (item *Item, err error) {
	if name == "" {
		err = errors.New("empty name")
		return
	}
	pipe, err := dao_pipe.GetPipeline(name)
	if err != nil {
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
func PipeStatus(name string) (status pipeline.Status, err error) {
	pipe, err := dao_pipe.GetPipeline(name)
	if err != nil {
		return
	}
	status = pipe.Status
	return
}
