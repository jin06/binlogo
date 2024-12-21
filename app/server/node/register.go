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
	closed    chan struct{}
	closeOnce sync.Once
}

func NewRegister(n *node.Node) *Register {
	r := &Register{
		node:   n,
		closed: make(chan struct{}),
	}
	return r
}

func (r *Register) Run(ctx context.Context) error {
	defer r.close()
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	ok, err := dao.RegisterNode(ctx, r.node)
	if !ok {
		return err
	}
	check := time.Now()
	for i := 0; i < 3; {
		select {
		case <-ctx.Done():
			return nil
		case <-r.closed:
			return nil
		case <-ticker.C:
			if time.Since(check) > time.Second*5 {
				return nil
			}
			err := dao.LeaseNode(ctx, r.node)
			if err == nil {
				check = time.Now()
			}
		}
	}
	return nil
}

func (r *Register) close() error {
	r.closeOnce.Do(func() {
		close(r.closed)
	})
	return nil
}
