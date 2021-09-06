package options

import "time"

type Options struct {
	Endpoints []string
	Timeout time.Duration
	Prefix string
}

type Option func(*Options)

func Endpoints(endpoints []string) Option {
	return func(options *Options) {
		options.Endpoints = endpoints
	}
}

func Timeout(timeout time.Duration) Option {
	return func(options *Options) {
		options.Timeout = timeout
	}
}

func Prefix(prefix string ) Option {
	return func(options *Options) {
		options.Prefix = prefix
	}
}
