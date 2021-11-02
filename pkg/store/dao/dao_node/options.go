package dao_node

import (
	"github.com/jin06/binlogo/pkg/node/role"
	"net"
)

type StatusReady bool

type Options struct {
	Key         string
	NodeName    string
	NodeIP      net.IP
	NodeVersion string
	NodeRole    role.Role
	StatusReady *boolVal
}

func WithKey(key string) Option {
	return func(opts *Options) {
		opts.Key = key
		return
	}
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
func WithNodeName(name string) Option {
	return func(options *Options) {
		options.NodeName = name
		return
	}
}

func WithNodeIP(ip net.IP) Option {
	return func(options *Options) {
		options.NodeIP = ip
	}
}

func WithNodeVersion(ver string) Option {
	return func(options *Options) {
		options.NodeVersion = ver
	}
}
