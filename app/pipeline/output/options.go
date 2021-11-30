package output

import (
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
)

// Options for configure output
type Options struct {
	Output       *pipeline.Output
	PipelineName string
}

// Option is a function for configure Options
type Option func(options *Options)

// OptionOutput sets Output
func OptionOutput(val *pipeline.Output) Option {
	return func(options *Options) {
		options.Output = val
	}
}

// OptionPipeName sets pipeline name
func OptionPipeName(name string) Option {
	return func(options *Options) {
		options.PipelineName = name
	}
}
