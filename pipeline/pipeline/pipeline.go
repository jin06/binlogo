package pipeline

import (
	"errors"
	"github.com/jin06/binlogo/pipeline/control"
	"github.com/jin06/binlogo/pipeline/filter"
	"github.com/jin06/binlogo/pipeline/input"
	"github.com/jin06/binlogo/pipeline/message"
	"github.com/jin06/binlogo/pipeline/output"
)

type Pipeline struct {
	Input       *input.Input
	Output      *output.Output
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

func (p *Pipeline) Init() (err error) {
	if p.Options.Pipeline == nil {
		return errors.New("lack pipeline info")
	}
	if p.Options.Mysql == nil {
		return errors.New("lack mysql info")
	}
	in := input.NewInput(
		input.OptionMysql(p.Options.Mysql),
		input.OptionPosition(p.Options.Position),
	)
	p.Input = in
	return
}

func NewPipeline(opt ...Option) (p *Pipeline) {
	options := Options{}
	for _, v := range opt {
		v(&options)
	}
	p = &Pipeline{
		Options: options,
	}
	return
}
