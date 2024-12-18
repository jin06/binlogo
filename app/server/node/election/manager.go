package election

import (
	"context"
	"sync"

	"github.com/jin06/binlogo/v2/internal/constant"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
)

// Manager is election cycle manager
type Manager struct {
	optNode   *node.Node
	roleCh    chan constant.Role
	closeOnce sync.Once
	stopped   chan struct{}
}

// NewManager returns a new manager
func NewManager(optNode *node.Node) (m *Manager) {
	m = &Manager{
		optNode: optNode,
		roleCh:  make(chan constant.Role, 1000),
		stopped: make(chan struct{}),
	}
	return
}

// Run election cycle
func (m *Manager) Run(ctx context.Context) {
	defer m.close()
	stx, cancel := context.WithCancel(ctx)
	defer cancel()
	for {
		// election := New(
		// 	OptionNode(m.optNode),
		// 	OptionTTL(5),
		// )
		election := NewElection(m.optNode)
		election.Run(stx, m.roleCh)
	}
}

// RoleCh return role chan
func (m *Manager) RoleCh() chan constant.Role {
	return m.roleCh
}

func (m *Manager) close() error {
	m.closeOnce.Do(func() {
		close(m.roleCh)
		close(m.stopped)
	})
	return nil
}

func (m *Manager) Stopped() chan struct{} {
	return m.stopped
}
