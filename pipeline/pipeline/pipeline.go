package pipeline

import (
	"github.com/jin06/binlogo/pipeline/control"
	"github.com/jin06/binlogo/pipeline/filter"
	"github.com/jin06/binlogo/pipeline/input"
	"github.com/jin06/binlogo/pipeline/message"
	"github.com/jin06/binlogo/pipeline/output"
)

type Pipeline struct {
	Input       *input.Controller
	Output      *output.Controller
	Filter      *filter.Controller
	DataLine    *DataLine
	ControlLine *ControlLine
	Options     Options
}

type DataLine struct {
	Input  chan *message.Message
	Filter chan *message.Message
	Output chan *message.Message
}

type ControlLine struct {
	Command chan *control.Command
}

func (p *Pipeline) Init() {

}
