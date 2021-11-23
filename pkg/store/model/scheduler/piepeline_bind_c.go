package scheduler

import (
	"github.com/jin06/binlogo/pkg/store/model"
)

// PipelineBindH deprecated
type PipelineBindH struct {
	*PipelineBind
	*model.Header
}

// NewPipelineBindH deprecated
func NewPipelineBindH() *PipelineBindH {
	r := &PipelineBindH{}
	r.PipelineBind = &PipelineBind{
		Bindings: map[string]string{},
	}
	r.Header = &model.Header{}
	return r
}
