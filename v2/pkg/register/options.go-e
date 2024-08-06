package register

// Option function to configure *Register
type Option func(r *Register)

// WithTTL sets ttl
func WithTTL(ttl int64) Option {
	return func(r *Register) {
		r.ttl = ttl
	}
}

// WithKey sets registerKey
func WithKey(key string) Option {
	return func(r *Register) {
		r.registerKey = key
	}
}

// WithData sets registerData
func WithData(data interface{}) Option {
	return func(r *Register) {
		r.registerData = data
	}
}
