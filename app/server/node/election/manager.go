package election

import (
	"context"
	"github.com/jin06/binlogo/pkg/node/role"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"sync"
)

// Manager is election cycle manager
type Manager struct {
	ctx      context.Context
	mutex    sync.Mutex
	optNode  *node.Node
	election *Election
}

// NewManager returns a new manager
func NewManager() (m *Manager) {
	m = &Manager{
		ctx:   context.Background(),
		mutex: sync.Mutex{},
	}
	return
}

// Run election cycle
func (m *Manager) Run(ctx context.Context) (err error) {
	myCtx, cancel := context.WithCancel(ctx)
	m.ctx = myCtx
	go func() {
		defer func() {
			cancel()
		}()
		m.election = New(
			OptionNode(m.optNode),
			OptionTTL(5),
		)
		m.election.Run(myCtx)
		select {
		case <-ctx.Done():
			{
				return
			}
		case <-m.Context().Done():
			{
				break
			}
		}
	}()
	return
}

// Context returns context
func (m *Manager) Context() context.Context {
	return m.ctx
}

// Role returns role
func (m *Manager) Role() role.Role {
	if m.election == nil {
		return role.FOLLOWER
	}
	return m.election.Role()
}
