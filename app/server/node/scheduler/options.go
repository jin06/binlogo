package scheduler

import (
	"github.com/jin06/binlogo/pkg/store/model/node"
)

// Options for configure scheduler
type Options struct {
	Node *node.Node
}

// Option is function thant sets Options
type Option func(opts *Options)
