package pipeline

import "github.com/jin06/binlogo/store/model"

type Options struct {
	Pipeline *model.Pipeline
	Mysql    *model.Mysql
	Position *model.Position
	Filters  []*model.Filter
}

type Option func(*Options)

func OptionPipeline(option *model.Pipeline) Option {
	return func(options *Options) {
		options.Pipeline = option
	}
}

func OptionMysql(option *model.Mysql) Option {
	return func(options *Options) {
		options.Mysql = option
	}
}

func OptionPosition(option *model.Position) Option {
	return func(options *Options) {
		options.Position = option
	}
}
