package mutex

import (
	"context"

	"github.com/jin06/binlogo/pkg/etcdclient"

	//"github.com/jin06/binlogo/pkg/store/etcd"
	"time"

	"go.etcd.io/etcd/client/v3/concurrency"
)

// Mutex Distributed lock encapsulating etcd
type Mutex struct {
	mutex   *concurrency.Mutex
	session *concurrency.Session
	timeout time.Duration
}

// New returns a new *Mutex
func New(key string, opts ...Option) (m *Mutex, err error) {
	m = &Mutex{
		timeout: time.Second * 5,
	}
	m.initOptions(opts...)
	cli, err := etcdclient.New()
	if err != nil {
		return
	}
	m.session, err = concurrency.NewSession(cli)
	if err != nil {
		return
	}
	m.mutex = concurrency.NewMutex(m.session, key)
	return
}

func (m *Mutex) initOptions(opts ...Option) {
	for _, v := range opts {
		v(m)
	}
}

// Lock lock
func (m *Mutex) Lock() (err error) {
	ctx, _ := context.WithTimeout(context.TODO(), m.timeout)
	err = m.mutex.Lock(ctx)
	return
}

// Unlock unlock
func (m *Mutex) Unlock() (err error) {
	err = m.mutex.Unlock(context.TODO())
	if err != nil {
		return
	}
	err = m.session.Close()
	return
}
