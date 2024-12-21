package election

import (
	"context"
	"sync"
	"time"

	"github.com/jin06/binlogo/v2/internal/constant"
	"github.com/jin06/binlogo/v2/pkg/store/dao"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
)

// Manager is election cycle manager
type Manager struct {
	optNode   *node.Node
	roleCh    chan constant.Role
	role      constant.Role
	mu        sync.Mutex
	closeOnce sync.Once
	closed    chan struct{}
}

// NewManager returns a new manager
func NewManager(optNode *node.Node) (m *Manager) {
	m = &Manager{
		optNode: optNode,
		roleCh:  make(chan constant.Role, 1000),
		closed:  make(chan struct{}),
	}
	return
}

// Run election cycle
func (m *Manager) Run(ctx context.Context) {
	defer m.close()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if m.role == constant.LEADER {
				key, err := dao.GetMasterLock(ctx)
				if err != nil {
					m.setRole(constant.FOLLOWER)
				}
				if key == m.optNode.Name {
					err := dao.LeaseMasterLock(ctx)
					if err != nil {
						m.setRole(constant.FOLLOWER)
					}
				}
			}
			if m.role == constant.FOLLOWER {
				err := dao.AcquireMasterLock(ctx, m.optNode)
				if err == nil {
					m.setRole(constant.LEADER)
				} else {
					m.setRole(constant.FOLLOWER)
				}
			}
		}
	}
}

// RoleCh return role chan
func (m *Manager) RoleCh() chan constant.Role {
	return m.roleCh
}

func (m *Manager) setRole(r constant.Role) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.role = r
	m.roleCh <- r
}

func (m *Manager) close() error {
	m.closeOnce.Do(func() {
		close(m.roleCh)
		close(m.closed)
	})
	return nil
}

func (m *Manager) Closed() chan struct{} {
	return m.closed
}
