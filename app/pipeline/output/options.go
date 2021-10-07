package output

import (
	model2 "github.com/jin06/binlogo/pkg/store/model"
)

type Options struct {
	Output *model2.Output
}

type Option func(options *Options)

func OptionOutput(val *model2.Output) Option {
	return func(options *Options) {
		options.Output = val
	}
}
