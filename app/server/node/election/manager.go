package election

import (
	"context"
	"github.com/jin06/binlogo/pkg/node/role"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"sync"
)

// Manager is election cycle manager
type Manager struct {
	optNode   *node.Node
	roleCh    chan role.Role
	cancel    context.CancelFunc
	closeOnce sync.Once
	stopped   chan struct{}
}

// NewManager returns a new manager
func NewManager(optNode *node.Node) (m *Manager) {
	m = &Manager{
		optNode: optNode,
		roleCh:  make(chan role.Role, 1000),
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
		election := New(
			OptionNode(m.optNode),
			OptionTTL(5),
		)

		go election.Run(stx)
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			case <-election.stopped:
				{
					break
				}
			case r, ok := <-election.RoleCh:
				{
					if !ok {
						break
					}

					m.roleCh <- r
				}
			}
		}
		election.close()
	}
}

// RoleCh return role chan
func (m *Manager) RoleCh() chan role.Role {
	return m.roleCh
}

func (m *Manager) close() error {
	m.closeOnce.Do(func() {
		close(m.stopped)
	})
	return nil
}

func (m *Manager) Stopped() chan struct{} {
	return m.stopped
}
