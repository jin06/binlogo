package node

import "context"

type NodeWatcher struct {
	Options *Options
}

func (w *NodeWatcher) Run(ctx context.Context) {
	return
}

func New() (w *NodeWatcher, err error) {
	return
}
