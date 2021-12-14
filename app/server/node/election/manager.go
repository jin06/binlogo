package election

import (
	"context"
	"sync"

	"github.com/jin06/binlogo/pkg/node/role"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/sirupsen/logrus"
)

// Manager is election cycle manager
type Manager struct {
	ctx      context.Context
	mutex    sync.Mutex
	optNode  *node.Node
	election *Election
	roleCh   chan role.Role
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
	myCtx, cancel := context.WithCancel(ctx)
	m.ctx = myCtx
	go func() {
		defer func() {
			defer func() {
				r := recover()
				if r != nil {
					logrus.Errorln("election manager panic, ", r)
				}
			}()
			cancel()
		}()
		for {
			en := New(
				OptionNode(m.optNode),
				OptionTTL(5),
			)

			m.election = en
			cCtx, _ := context.WithCancel(myCtx)
			en.Run(cCtx)
			go func() {
				defer func() {
					if r := recover(); r != nil {
						logrus.Errorln("election manager panic, ", r)
					}
				}()
				for {
					select {
					case <-myCtx.Done():
						{
							return
						}
					case <-en.Context().Done():
						{
							return
						}
					case r, ok := <-en.RoleCh:
						{
							if !ok {
								return
							}
							if r == m.election.Role() {
								m.roleCh <- r
							}
						}
					}
				}
			}()
			select {
			case <-myCtx.Done():
				{
					return
				}
			case <-ctx.Done():
				{
					return
				}
			case <-en.Context().Done():
				{
					break
				}
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

// RoleCh return role chan
func (m *Manager) RoleCh() chan role.Role {
	return m.roleCh
}
