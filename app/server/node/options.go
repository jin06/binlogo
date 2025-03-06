package node

import (
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
	"github.com/sirupsen/logrus"
)

// Options Node's options
type Options struct {
	Node *node.Node
	Log  *logrus.Entry
}

// Option function config Options
type Option func(options *Options)
