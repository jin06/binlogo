package node
type Options struct {
	Etcd *struct {
		Endpoints []string
	}
}

type Option func(opts *Options)

func WithEtcdEndpoints(endpoints []string) Option {
	return func(opts *Options) {
		if opts.Etcd == nil {
			opts.Etcd = &struct{ Endpoints []string }{Endpoints: []string{}}
			opts.Etcd.Endpoints = endpoints
		}
	}
}
