package node

import (
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
)

// Options Node's options
type Options struct {
	Node *node.Node
}

// Option function config Options
type Option func(options *Options)

// OptionNode sets Options Node
func OptionNode(node *node.Node) Option {
	return func(options *Options) {
		options.Node = node
	}
}
