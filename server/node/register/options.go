package register

import (
	"github.com/jin06/binlogo/store/etcd"
	"github.com/jin06/binlogo/store/model"
	"time"
)

type Option func(options *Register)

func OptionNode(node *model.Node) Option {
	return func(r *Register) {
		r.node = node
	}
}

func OptionTTL(ttl int64) Option {
	return func(r *Register) {
		r.ttl = ttl
	}
}

func OptionETCD(etcd *etcd.ETCD) Option {
	return func(r *Register) {
		r.etcd = etcd
	}
}

func OptionLeaseDuration(t time.Duration) Option {
	return func (r *Register) {
		r.leaseDuration = t
	}
}
