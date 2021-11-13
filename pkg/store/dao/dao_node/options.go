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

func WithNodeName(name string) Option {
	return func(options *Options) {
		options.NodeName = name
		return
	}
}

//func WithNodeIP(ip net.IP) Option {
//	return func(options *Options) {
//		options.NodeIP = ip
//	}
//}
//
//func WithNodeVersion(ver string) Option {
//	return func(options *Options) {
//		options.NodeVersion = ver
//	}
//}
