package dao

type Options struct {
	Key string
}
type Option func (*Options)

func GetOptions(args ...Option) *Options{
	opts := &Options{}
	for _,v := range args {
		v(opts)
	}
	return opts
}

func WithKey(key string) Option {
	return func(opts *Options) {
		opts.Key = key
		return
	}
}


