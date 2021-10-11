package register

import (
	"context"
	"errors"
	"fmt"
	etcd2 "github.com/jin06/binlogo/pkg/store/etcd"
	model2 "github.com/jin06/binlogo/pkg/store/model"
	"github.com/jin06/binlogo/pkg/util/ip"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
	"os"
	"time"
)

const ttl = 5

func New(opts ...Option) (r *Register) {
	r = &Register{
		ttl:           ttl,
		etcd:          etcd2.E,
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
	node          *model2.Node
	ttl           int64
	etcd          *etcd2.ETCD
}

func (r *Register) Run(ctx context.Context) (err error) {
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
			case <-ctx.Done():
				{
					logrus.Warn("register exit")
					return
				}
			}
		}
	}()

	return
}

func (r *Register) reg() (err error) {

	if r.node.Ip == nil {
		r.node.Ip, err = ip.LocalIp()
		if err != nil {
			logrus.Error("Get local ip error ", err)
			os.Exit(1)
			return
		}
	}

	r.lease = clientv3.NewLease(r.etcd.Client)
	rep, err := r.lease.Grant(context.TODO(), r.ttl)
	if err != nil {
		return
	}
	r.leaseID = rep.ID
	ok, err := r.etcd.Create(r.node, clientv3.WithLease(r.leaseID))
	logrus.Debug("register tick")
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
