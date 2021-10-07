package register

import (
	"context"
	"errors"
	"fmt"
	"github.com/jin06/binlogo/store/etcd"
	"github.com/jin06/binlogo/store/model"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
	"time"
)

const ttl = 5

func New(opts ...Option) (r *Register) {
	r = &Register{
		ttl:           ttl,
		etcd:          etcd.E,
		leaseDuration: time.Second,
	}
	logrus.Debug("Node register lease duration ", r.leaseDuration)
	for _, v := range opts {
		v(r)
	}
	logrus.Debug("Node register lease duration ", r.leaseDuration)
	return
}

type Register struct {
	lease         clientv3.Lease
	leaseID       clientv3.LeaseID
	leaseDuration time.Duration
	node          *model.Node
	ttl           int64
	etcd          *etcd.ETCD
}

func (r *Register) Run() (err error) {
	err = r.reg()
	if err != nil {
		return
	}
	go func() {
		tc := time.NewTicker(r.leaseDuration)
		for {
			select {
			case <-tc.C:
				{
					_ = r.keep()
				}
			}
		}
	}()

	return
}

func (r *Register) reg() (err error) {
	r.lease = clientv3.NewLease(r.etcd.Client)
	rep, err := r.lease.Grant(context.TODO(), r.ttl)
	if err != nil {
		return
	}
	r.leaseID = rep.ID
	ok, err := r.etcd.Create(r.node, clientv3.WithLease(r.leaseID))
	fmt.Println("tick")
	if err != nil {
		fmt.Println("reg error: ", err)
		return
	}
	if !ok {
		return errors.New("register failed")
	}
	return
}

func (r *Register) keep() (err error) {
	resp, err := r.lease.KeepAliveOnce(context.TODO(), r.leaseID)
	if err != nil {
		logrus.Error(err)
	}
	logrus.Debugf("Node lease %v seconds.", resp.TTL)
	return
}
func (r *Register) revoke() (err error) {
	_, err = r.etcd.Client.Revoke(context.TODO(), r.leaseID)
	return
}
