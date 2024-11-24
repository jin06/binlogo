package node

import (
	"context"
	"sync"
	"time"

	"github.com/jin06/binlogo/v2/pkg/store/dao"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
)

type Register struct {
	node      *node.Node
	stopped   chan struct{}
	closeOnce sync.Once
}

func NewRegister(n *node.Node) *Register {
	r := &Register{
		node:    n,
		stopped: make(chan struct{}),
	}
	return r
}

func (r *Register) Run(ctx context.Context) error {
	defer r.close()
	ticker := time.NewTicker(time.Second * 1)
	ok, err := dao.RegisterNode(ctx, r.node)
	if !ok {
		return err
	}
	for i := 0; i < 3; {
		select {
		case <-ctx.Done():
			return nil
		case <-r.stopped:
			return nil
		case <-ticker.C:
			exist, err := dao.LeaseNode(ctx, r.node)
			if err != nil {
				i++
			} else {
				if exist {
					i = 0
				} else {
					return nil
				}
			}
		}
	}
	return nil
}

func (r *Register) close() error {
	r.closeOnce.Do(func() {
		close(r.stopped)
	})
	return nil
}
