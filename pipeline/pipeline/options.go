package pipeline

import "github.com/jin06/binlogo/store/model"

func NewPipeline(opt ...Option) (p *Pipeline) {
	options := Options{}
	for _, v := range opt {
		v(&options)
	}
	p = &Pipeline{
		Options: options,
	}
	return
}

type Options struct {
	Pipeline *model.Pipeline
	Mysql    *model.Mysql
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
