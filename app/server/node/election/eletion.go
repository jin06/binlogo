package election

import (
	"context"
	"sync"
	"time"

	"github.com/jin06/binlogo/v2/internal/constant"
	"github.com/jin06/binlogo/v2/pkg/store/dao"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
	"github.com/sirupsen/logrus"
)

type Election struct {
	node      *node.Node
	mu        sync.Mutex
	role      constant.Role
	roleCh    chan constant.Role
	ttl       time.Duration
	closeOnce sync.Once
	closing   chan struct{}
	closed    chan struct{}
}

func NewElection(node *node.Node) *Election {
	e := &Election{
		node:    node,
		ttl:     5,
		roleCh:  make(chan constant.Role, 100),
		closing: make(chan struct{}),
		closed:  make(chan struct{}),
	}
	return e
}

func (e *Election) setRole(r constant.Role) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.role = r
	e.roleCh <- r
}

func (e *Election) Run(ctx context.Context, ch chan constant.Role) error {
	defer close(e.closed)
	defer close(e.roleCh)
	e.roleCh = ch
	return nil
}

func (e *Election) campaign(ctx context.Context) {
	logrus.Info("Election Cycle: ")
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-e.closing:
			return
		case <-ticker.C:
			err := dao.AcquireMasterLock(ctx, e.node)
			if err == nil {
				e.setRole(constant.LEADER)
			} else {
				e.setRole(constant.FOLLOWER)
			}
		}
	}
}

func (e *Election) Close() (err error) {
	e.closeOnce.Do(func() {
		e.setRole(constant.FOLLOWER)
		close(e.closing)
	})
	return
}

func (e *Election) Closed() chan struct{} {
	return e.closed
}
