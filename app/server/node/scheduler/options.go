package scheduler

import (
	"github.com/jin06/binlogo/pkg/store/model/node"
)

type Options struct {
	Node *node.Node
}

type Option func(opts *Options)

func WithNode(node *node.Node) Option {
	return func(opts *Options) {
		opts.Node = node
	}
}
