package register

import (
	"github.com/jin06/binlogo/pkg/store/model/node"
)

type Option func(options *Register)

func OptionNode(node *node.Node) Option {
	return func(r *Register) {
		r.node = node
	}
}


