package register

import (
	"github.com/jin06/binlogo/pkg/store/model/node"
	"time"
)

type Option func(options *Register)

func OptionNode(node *node.Node) Option {
	return func(r *Register) {
		r.node = node
	}
}

func OptionLeaseDuration(t time.Duration) Option {
	return func (r *Register) {
		r.leaseDuration = t
	}
}
