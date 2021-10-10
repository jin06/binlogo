package scheduler

import "github.com/jin06/binlogo/pkg/store/model"

type Options struct {
	Node *model.Node
}

type Option func(opts *Options)

func WithNode(node *model.Node) Option {
	return func(opts *Options) {
		opts.Node = node
	}
}
