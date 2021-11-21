package manager

import (
	"context"
	"sync"
)

type Manager interface {
	Start(ctx context.Context) chan error
}

type Status byte

const (
	STOP  Status = '1'
	START Status = '2'
)

type Base struct {
	Status Status
	mutex sync.Mutex
}

func (base *Base) Lock() {
	base.mutex.Lock()
}
func (base *Base) Unlock() {
	base.mutex.Unlock()
}
