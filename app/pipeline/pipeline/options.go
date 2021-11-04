package pipeline

import (
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
)

type Options struct {
	Pipeline *pipeline.Pipeline
	Position *pipeline.Position
}

type Option func(*Options)

func OptionPipeline(option *pipeline.Pipeline) Option {
	return func(options *Options) {
		options.Pipeline = option
	}
}

func OptionPosition(option *pipeline.Position) Option {
	return func(options *Options) {
		options.Position = option
	}
}

