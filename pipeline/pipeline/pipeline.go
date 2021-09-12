package pipeline

import (
	"github.com/jin06/binlogo/pipeline/filter"
	"github.com/jin06/binlogo/pipeline/input"
	"github.com/jin06/binlogo/pipeline/output"
)

type Pipeline struct {
	Input  *input.Manage
	Output *output.Manage
	Filter *filter.Manage
	Options
}

type Options struct {
	ID string
}

type Option func(*Options)

func OptionID(id string) Option {
	return func(ops *Options) {
		ops.ID = id
	}
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
