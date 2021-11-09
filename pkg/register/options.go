package register

type Option func(r *Register)

func WithTTL(ttl int64) Option {
	return func(r *Register) {
		r.ttl = ttl
	}
}

func WithKey(key string) Option {
	return func(r *Register) {
		r.registerKey = key
	}
}

func WithData(data interface{}) Option {
	return func(r *Register) {
		r.registerData = data
	}
}
