package mutex

import (
	"context"
	"github.com/jin06/binlogo/pkg/etcd_client"

	//"github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/coreos/etcd/clientv3/concurrency"
	"time"
)

type Mutex struct {
	mutex   *concurrency.Mutex
	session *concurrency.Session
	timeout time.Duration
}

func New(key string, opts ...Option) (m *Mutex, err error) {
	m = &Mutex{
		timeout: time.Second * 5,
	}
	m.initOptions(opts...)
	cli, err := etcd_client.New()
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

func (m *Mutex) Lock() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	err = m.mutex.Lock(ctx)
	return
}

func (m *Mutex) Unlock() (err error) {
	err = m.mutex.Unlock(context.Background())
	if err != nil {
		return
	}
	err = m.session.Close()
	return
}
