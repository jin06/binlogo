package manager_event

import (
	"sync"

	"github.com/jin06/binlogo/v2/app/server/node/manager"
)

// Manager is event manager
type Manager struct {
	*manager.Base
	Exit     bool
	stopOnce sync.Once
	stopping chan struct{}
}

// New returns a new Manager
func New() (m *Manager) {
	m = &Manager{
		Base:     &manager.Base{Status: manager.STOP},
		Exit:     false,
		stopping: make(chan struct{}),
	}
	return
}

// Stop stop working when lost leader
func (m *Manager) Stop() {
	m.stop()
}

func (m *Manager) stop() {
	m.stopOnce.Do(func() {
		close(m.stopping)
	})
}
