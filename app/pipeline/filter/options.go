package filter

import (
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
)

// Options is an interface abstraction for dynamic configuration
type Options struct {
	Pipe *pipeline.Pipeline
}

// Option function config Options
type Option func(options *Options)

// WithPipe sets pipeline to Options
func WithPipe(p *pipeline.Pipeline) Option {
	return func(options *Options) {
		options.Pipe = p
	}
}
