package output

import "github.com/jin06/binlogo/store/model"

type Options struct {
	Output *model.Output
}

type Option func(options *Options)

func OptionOutput(val *model.Output) Option {
	return func(options *Options) {
		options.Output = val
	}
}
