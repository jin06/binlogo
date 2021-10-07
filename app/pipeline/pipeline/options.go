package pipeline

import "github.com/jin06/binlogo/store/model"

type Options struct {
	Pipeline *model.Pipeline
	Position *model.Position
}

type Option func(*Options)

func OptionPipeline(option *model.Pipeline) Option {
	return func(options *Options) {
		options.Pipeline = option
	}
}

func OptionPosition(option *model.Position) Option {
	return func(options *Options) {
		options.Position = option
	}
}

