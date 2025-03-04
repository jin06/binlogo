package output

import (
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
)

// Options for configure output
type Options struct {
	Output       *pipeline.Output
	PipelineName string
	MysqlMode    pipeline.Mode
}

// Option is a function for configure Options
type Option func(options *Options)

// OptionOutput sets Output
func OptionOutput(val pipeline.Output) Option {
	return func(options *Options) {
		options.Output = &val
	}
}

// OptionPipeName sets pipeline name
func OptionPipeName(name string) Option {
	return func(options *Options) {
		options.PipelineName = name
	}
}

// OptionPipeMode sets mysql mode
func OptionMysqlMode(mode pipeline.Mode) Option {
	return func(options *Options) {
		options.MysqlMode = mode
	}
}
