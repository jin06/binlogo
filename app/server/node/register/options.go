package register

import (
	etcd2 "github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"time"
)

type Option func(options *Register)

func OptionNode(node *node.Node) Option {
	return func(r *Register) {
		r.node = node
	}
}

func OptionTTL(ttl int64) Option {
	return func(r *Register) {
		r.ttl = ttl
	}
}

func OptionETCD(etcd *etcd2.ETCD) Option {
	return func(r *Register) {
		r.etcd = etcd
	}
}

func OptionLeaseDuration(t time.Duration) Option {
	return func (r *Register) {
		r.leaseDuration = t
	}
}
