package pipeline

import (
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
	"github.com/sirupsen/logrus"
)

// Options for pipeline
type Options struct {
	Pipeline *pipeline.Pipeline
	Position *pipeline.Position
	Log      *logrus.Entry
}

// Option is a function for configure Options
type Option func(*Options)

func OptionLog(option *logrus.Entry) Option {
	return func(options *Options) {
		options.Log = option
	}
}

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
