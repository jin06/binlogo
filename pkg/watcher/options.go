package watcher

type Options struct {
	Key string
}

type Option func(*Options)

func NewOptions(opt ...Option) *Options {
	opts := &Options{}
	for _, o := range opt {
		o(opts)
	}
	return opts
}

func WithKey(k string) Option {
	return func(options *Options) {
		options.Key = k
	}
}
