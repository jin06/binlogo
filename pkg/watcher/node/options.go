package node

// Options of node watcher
type Options struct {
	Etcd *struct {
		Endpoints []string
	}
}

// Option function of configure Options
type Option func(opts *Options)

// WithEtcdEndpoints sets Endpoints of Options.Etcd.Endpoints
func WithEtcdEndpoints(endpoints []string) Option {
	return func(opts *Options) {
		if opts.Etcd == nil {
			opts.Etcd = &struct{ Endpoints []string }{Endpoints: []string{}}
			opts.Etcd.Endpoints = endpoints
		}
	}
}
