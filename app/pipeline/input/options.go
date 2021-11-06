package input

import (
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
)

type Options struct {
	Position *pipeline.Position
	Mysql    *pipeline.Mysql
	Pipeline *pipeline.Pipeline
}

type Option func(options *Options)

func OptionMysql(mysql *pipeline.Mysql) Option {
	return func(options *Options) {
		options.Mysql = mysql
	}
}

func OptionPosition(position *pipeline.Position) Option {
	return func(options *Options) {
		options.Position = position
	}
}

func OptionsPipeline(p *pipeline.Pipeline) Option {
	return func(options *Options) {
		options.Pipeline = p
	}
}
