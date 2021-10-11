package pipeline

import (
	model2 "github.com/jin06/binlogo/pkg/store/model"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
)

type Options struct {
	Pipeline *pipeline.Pipeline
	Position *model2.Position
}

type Option func(*Options)

func OptionPipeline(option *pipeline.Pipeline) Option {
	return func(options *Options) {
		options.Pipeline = option
	}
}

func OptionPosition(option *model2.Position) Option {
	return func(options *Options) {
		options.Position = option
	}
}

