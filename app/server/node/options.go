package node

import (
	model2 "github.com/jin06/binlogo/pkg/store/model"
)

type Options struct {
	Node *model2.Node
}

type Option func(options *Options)

func OptionNode(node *model2.Node) Option {
	return func(options *Options) {
		options.Node = node
	}
}



