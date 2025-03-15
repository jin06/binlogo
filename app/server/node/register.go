package node

import (
	"context"
	"sync"
	"time"

	"github.com/jin06/binlogo/v2/pkg/store/dao"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
	"github.com/sirupsen/logrus"
)

type Register struct {
	node         *node.Node
	closed       chan struct{}
	completeOnce sync.Once
	closing      chan struct{}
	closeOnce    sync.Once
	isClosed     bool
	log          *logrus.Entry
}

func NewRegister(log *logrus.Entry, n *node.Node) *Register {
	r := &Register{
		node:   n,
		closed: make(chan struct{}),
	}
	r.log = log
	// r.log = log.WithField("NodeRegSeq", )
	return r
}

func (r *Register) Run(ctx context.Context) error {
	defer r.CompleteClose()
	defer r.Close()
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	ok, err := dao.RegisterNode(ctx, r.node)
	if !ok {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-r.closing:
			return nil
		case <-ticker.C:
			if err := dao.LeaseNode(ctx, r.node); err != nil {
				panic(err)
				// return err
			}
		}
	}
}

func (r *Register) Close() {
	r.closeOnce.Do(func() {
		close(r.closing)
	})
}

func (r *Register) Closed() chan struct{} {
	return r.closed
}

func (r *Register) CompleteClose() {
	r.completeOnce.Do(func() {
		r.isClosed = true
		close(r.closed)
	})
}
