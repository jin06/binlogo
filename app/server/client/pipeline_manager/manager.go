package pipeline_manager

import (
	"context"
	"sync"
)

type Manager struct {
	status      Status
	runMutex    sync.Mutex
	statusMutex sync.Mutex
}

type Status byte

const (
	STATUS_RUN  Status = 1
	STATUS_STOP Status = 2
)

func New() *Manager {
	m := &Manager{
		status:      STATUS_STOP,
		runMutex:    sync.Mutex{},
		statusMutex: sync.Mutex{},
	}
	return m
}

func (m *Manager) setStatus(status Status) {
	m.statusMutex.Lock()
	defer m.statusMutex.Unlock()
	m.status = status
}

func (m *Manager) Status() Status {
	return m.status
}

func (m *Manager) Run(ctx context.Context) (err error) {
	m.runMutex.Lock()
	defer m.runMutex.Unlock()
	return
}

func (m *Manager) Stop(ctx context.Context) (err error) {
	m.runMutex.Lock()
	defer m.runMutex.Unlock()
	return
}
