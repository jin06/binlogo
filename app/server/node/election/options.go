package election

import "github.com/jin06/binlogo/pkg/store/model"

type Option func(e *Election)

func OptionNode(node *model.Node) Option {
	return func(e *Election) {
		e.node = node
	}
}

func OptionTTL(ttl int) Option {
	return func(e *Election) {
		e.ttl = ttl
	}
}

func OptionPrefix(prefix string) Option {
	return func(e *Election) {
		e.prefix = prefix
	}
}

