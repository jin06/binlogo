package pipeline

import (
	"fmt"
	"github.com/jin06/binlogo/pkg/store/model"
)

type PipelineH struct {
	*Pipeline
	*model.Header `json:"header"`
}

//func (m *PipelineH) GetHeader() *model.Header {
//	return m.Header
//}
//
//func (m *PipelineH) SetHeader(h *model.Header) {
//	m.Header = h
//}

func NewPipelineH() *PipelineH {
	p := &PipelineH{
		&Pipeline{},
		&model.Header{},
	}
	return p
}

func (m *PipelineH) String() string{
	return fmt.Sprintf("pipeline : %v \nheader : %v \n", m.Pipeline, m.Header)
}
