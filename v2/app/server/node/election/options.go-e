package election

import (
	"github.com/jin06/binlogo/pkg/store/model/node"
)

// Option is function of set election
type Option func(e *Election)

// OptionNode sets election node
func OptionNode(node *node.Node) Option {
	return func(e *Election) {
		e.node = node
	}
}

// OptionTTL sets election ttl
func OptionTTL(ttl int) Option {
	return func(e *Election) {
		e.ttl = ttl
	}
}
