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

// Manager is election cycle manager
type Manager struct {
	optNode      *node.Node
	roleCh       chan constant.Role
	role         constant.Role
	mu           sync.Mutex
	closing      chan struct{}
	closeOnce    sync.Once
	closed       chan struct{}
	completeOnce sync.Once
	log          *logrus.Entry
}

// NewManager returns a new manager
func NewManager(ctx context.Context, optNode *node.Node) (m *Manager) {
	m = &Manager{
		optNode: optNode,
		roleCh:  make(chan constant.Role, 1000),
		closing: make(chan struct{}),
		closed:  make(chan struct{}),
	}
	newCtx := context.WithValue(ctx, "Election CreateTime", time.Now())
	m.log = logrus.WithContext(newCtx)
	return
}

// Run election cycle
func (m *Manager) Run(ctx context.Context) {
	defer m.CompleteClose()
	defer m.Close()
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
			} else {
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

func (m *Manager) GetRole() constant.Role {
	return m.role
}

func (m *Manager) setRole(r constant.Role) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.role = r
	m.roleCh <- r
}

func (m *Manager) CompleteClose() {
	m.completeOnce.Do(func() {
		close(m.closed)
	})
}

func (m *Manager) Close() error {
	m.closeOnce.Do(func() {
		close(m.roleCh)
		close(m.closing)
	})
	return nil
}

func (m *Manager) Closed() chan struct{} {
	return m.closed
}
