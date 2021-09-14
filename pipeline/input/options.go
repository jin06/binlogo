package input

import "github.com/jin06/binlogo/store/model"

type Options struct {
	Position *model.Position
	Mysql    *model.Mysql
}

type Option func(options *Options)

func OptionMysql(mysql *model.Mysql) Option {
	return func(options *Options) {
		options.Mysql = mysql
	}
}

func OptionPosition(position *model.Position) Option {
	return func(options *Options) {
		options.Position = position
	}
}
