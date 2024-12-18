package election

import (
	"context"
	"sync"
	"time"

	"github.com/jin06/binlogo/v2/pkg/node/role"
	"github.com/jin06/binlogo/v2/pkg/store/dao"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
	"github.com/sirupsen/logrus"
)

type Election struct {
	node      *node.Node
	mu        sync.Mutex
	role      role.Role
	roleCh    chan role.Role
	ttl       time.Duration
	closeOnce sync.Once
	stopped   chan struct{}
}

func NewElection(node *node.Node) *Election {
	e := &Election{
		node:    node,
		ttl:     5,
		roleCh:  make(chan role.Role, 100),
		stopped: make(chan struct{}, 0),
	}
	return e
}

func (e *Election) setRole(r role.Role) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.role = r
	e.roleCh <- r
}

func (e *Election) Run(ctx context.Context, ch chan role.Role) error {
	defer e.close()
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
		case <-e.stopped:
			return
		case <-ticker.C:
			err := dao.AcquireMasterLock(ctx, e.node)
			if err == nil {
				e.setRole(role.LEADER)
			} else {
				e.setRole(role.FOLLOWER)
			}
		}
	}
}

func (e *Election) lease(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-ctx.Done():
		case <-e.stopped:
		case <-ticker.C:
		}
	}
}

func (e *Election) close() (err error) {
	e.closeOnce.Do(func() {
		e.setRole(role.FOLLOWER)
		close(e.stopped)
	})
	return
}
