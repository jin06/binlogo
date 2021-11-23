package manager

import (
	"context"
	"sync"
)

// Manager base interface
type Manager interface {
	Start(ctx context.Context) chan error
}

// Status of Manager
type Status byte

const (
	// Status stop
	STOP Status = '1'
	// Status start
	START Status = '2'
)

// Base a Manager implement
type Base struct {
	Status Status
	mutex  sync.Mutex
}

// Lock lock the Manger
func (base *Base) Lock() {
	base.mutex.Lock()
}

// Unlock unlock the Manager
func (base *Base) Unlock() {
	base.mutex.Unlock()
}
