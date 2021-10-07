package node

import "github.com/jin06/binlogo/store/model"

type Options struct {
	Node *model.Node
}

type Option func(options *Options)

func OptionNode(node *model.Node) Option {
	return func(options *Options) {
		options.Node = node
	}
}



