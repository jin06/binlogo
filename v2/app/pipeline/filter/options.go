package filter

import (
	pipeline2 "github.com/jin06/binlogo/pkg/store/model/pipeline"
)

// Options is an interface abstraction for dynamic configuration
type Options struct {
	Pipe *pipeline2.Pipeline
}

// Option function config Options
type Option func(options *Options)

// WithPipe sets pipeline to Options
func WithPipe(p *pipeline2.Pipeline) Option {
	return func(options *Options) {
		options.Pipe = p
	}
}
