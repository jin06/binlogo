package scheduler

import (
	"github.com/jin06/binlogo/pkg/store/model"
)

type PipelineBindH struct {
	*PipelineBind
	*model.Header
}

func NewPipelineBindH() *PipelineBindH {
	r := &PipelineBindH{}
	r.PipelineBind = &PipelineBind{
		Bindings: map[string]string{},
	}
	r.Header = &model.Header{}
	return r
}
