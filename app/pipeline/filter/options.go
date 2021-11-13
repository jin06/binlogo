package filter

import (
	pipeline2 "github.com/jin06/binlogo/pkg/store/model/pipeline"
)

type Options struct {
	Pipe *pipeline2.Pipeline
}

type Option func(options *Options)

func WithPipe(p *pipeline2.Pipeline) Option {
	return func(options *Options) {
		options.Pipe = p
	}
}
