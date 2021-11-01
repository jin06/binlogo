package dao_node

type StatusReady bool

type Options struct {
	StatusReady *boolVal
}

type boolVal struct {
	val bool
}

type Option func(*Options)

func GetOptions(args ...Option) *Options {
	opts := &Options{}
	for _, v := range args {
		v(opts)
	}
	return opts
}

func WithStatusReady(b bool) Option {
	return func(options *Options) {
		options.StatusReady = &boolVal{
			val: b,
		}
	}
}
