package output

import (
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
)

type Options struct {
	Output *pipeline.Output
}

type Option func(options *Options)

func OptionOutput(val *pipeline.Output) Option {
	return func(options *Options) {
		options.Output = val
	}
}
