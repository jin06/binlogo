package pipeline

import (
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
)

// Options for pipeline
type Options struct {
	Pipeline *pipeline.Pipeline
	Position *pipeline.Position
}

// Option is a function for configure Options
type Option func(*Options)

// OptionPipeline sets pipeline
func OptionPipeline(option *pipeline.Pipeline) Option {
	return func(options *Options) {
		options.Pipeline = option
	}
}

// OptionPosition sets position
func OptionPosition(option *pipeline.Position) Option {
	return func(options *Options) {
		options.Position = option
	}
}
