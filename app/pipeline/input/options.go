package input

import (
	model2 "github.com/jin06/binlogo/pkg/store/model"
)

type Options struct {
	Position *model2.Position
	Mysql    *model2.Mysql
}

type Option func(options *Options)

func OptionMysql(mysql *model2.Mysql) Option {
	return func(options *Options) {
		options.Mysql = mysql
	}
}

func OptionPosition(position *model2.Position) Option {
	return func(options *Options) {
		options.Position = position
	}
}
