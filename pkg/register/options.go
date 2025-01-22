package register

import (
	"time"

	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
)

// Option function to configure *Register
type Option func(r *Register)

// WithTTL sets ttl
func WithTTL(ttl time.Duration) Option {
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
func WithData(data *pipeline.Instance) Option {
	return func(r *Register) {
		r.registerData = data
	}
}
