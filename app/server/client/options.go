package client

type Options struct {

}

type Option func(options *Options)

func NewOptions(opts ...Option) *Options {
	options := &Options{}
	for _, v := range opts {
		v(options)
	}
	return options
}
