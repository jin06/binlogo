package election

import (
	"context"
	"github.com/jin06/binlogo/pkg/node/role"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/sirupsen/logrus"
	"sync"
)

// Manager is election cycle manager
type Manager struct {
	ctx       context.Context
	mutex     sync.Mutex
	optNode   *node.Node
	election  *Election
	roleCh    chan role.Role
	cancel    context.CancelFunc
	closeOnce sync.Once
}

// NewManager returns a new manager
func NewManager(optNode *node.Node) (m *Manager) {
	m = &Manager{
		ctx:     context.Background(),
		mutex:   sync.Mutex{},
		optNode: optNode,
		roleCh:  make(chan role.Role, 1000),
	}
	return
}

// Run election cycle
func (m *Manager) Run(ctx context.Context) {
	m.ctx, m.cancel = context.WithCancel(ctx)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorln("election manager panic, ", r)
				panic(r)
			}
			m.close()
		}()
		for {
			election := New(
				OptionNode(m.optNode),
				OptionTTL(5),
			)

			m.election = election
			election.Run(m.ctx)
			go func() {
				defer func() {
					if r := recover(); r != nil {
						logrus.Errorln("election manager panic, ", r)
					}
					election.close()
				}()
				for {
					select {
					case <-m.ctx.Done():
						{
							return
						}
					case <-election.context().Done():
						{
							return
						}
					case r, ok := <-election.RoleCh:
						{
							if !ok {
								return
							}
							if r == m.election.getRole() {
								m.roleCh <- r
							}
						}
					}
				}
			}()
			<-election.context().Done()
			break
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
	return m.election.getRole()
}

// RoleCh return role chan
func (m *Manager) RoleCh() chan role.Role {
	return m.roleCh
}

func (m *Manager) close() error {
	m.closeOnce.Do(func() {
		m.cancel()
		m.election.close()
	})
	return nil
}
