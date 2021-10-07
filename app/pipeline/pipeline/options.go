package pipeline

import (
	model2 "github.com/jin06/binlogo/pkg/store/model"
)

type Options struct {
	Pipeline *model2.Pipeline
	Position *model2.Position
}

type Option func(*Options)

func OptionPipeline(option *model2.Pipeline) Option {
	return func(options *Options) {
		options.Pipeline = option
	}
}

func OptionPosition(option *model2.Position) Option {
	return func(options *Options) {
		options.Position = option
	}
}

