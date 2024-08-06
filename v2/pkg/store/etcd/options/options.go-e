package options

import "time"

// Options for ETCD's Options
type Options struct {
	Endpoints []string
	Timeout   time.Duration
	Prefix    string
}

// Option is function set Options
type Option func(*Options)

// Endpoints sets Endpoints to Options
func Endpoints(endpoints []string) Option {
	return func(options *Options) {
		options.Endpoints = endpoints
	}
}

// Timeout sets Timeout to Options
func Timeout(timeout time.Duration) Option {
	return func(options *Options) {
		options.Timeout = timeout
	}
}

// Prefix sets Prefix to Options
func Prefix(prefix string) Option {
	return func(options *Options) {
		options.Prefix = prefix
	}
}
