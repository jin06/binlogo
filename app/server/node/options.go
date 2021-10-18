package node

import (
	"github.com/jin06/binlogo/pkg/store/model/node"
)

type Options struct {
	Node *node.Node
}

type Option func(options *Options)

func OptionNode(node *node.Node) Option {
	return func(options *Options) {
		options.Node = node
	}
}



