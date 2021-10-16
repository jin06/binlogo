package input

import (
	model2 "github.com/jin06/binlogo/pkg/store/model"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
)

type Options struct {
	Position *model2.Position
	Mysql    *pipeline.Mysql
}

type Option func(options *Options)

func OptionMysql(mysql *pipeline.Mysql) Option {
	return func(options *Options) {
		options.Mysql = mysql
	}
}

func OptionPosition(position *model2.Position) Option {
	return func(options *Options) {
		options.Position = position
	}
}
