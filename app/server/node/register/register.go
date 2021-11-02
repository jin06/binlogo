package register

import (
	"context"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/store/dao/dao_cluster"
	etcd2 "github.com/jin06/binlogo/pkg/store/etcd"
	node2 "github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/store/model/register"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
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
	node          *node2.Node
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
					er := r.keep()
					if er != nil {
						blog.Error(er)
					}
					if er == rpctypes.ErrLeaseNotFound {
						err = r.reg()
						blog.Error(err)
					}

					ok, err2 := r.watch()
					if err2 != nil {
						blog.Error(err2)
					}
					if !ok {
						err = r.reg()
					}
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
	r.lease = clientv3.NewLease(r.etcd.Client)
	if rep, err2 := r.lease.Grant(context.TODO(), r.ttl); err2 != nil {
		return err2
	} else {
		r.leaseID = rep.ID
	}
	err = dao_cluster.RegRNode(&register.RegisterNode{Name: r.node.Name}, r.leaseID)
	return
}

func (r *Register) keep() (err error) {
	resp, err := r.lease.KeepAliveOnce(context.TODO(), r.leaseID)
	if err != nil {
		//logrus.Error(err)
		return err
	}
	logrus.Debugf("Node lease %v seconds.", resp.TTL)
	return
}
func (r *Register) revoke() (err error) {
	_, err = r.etcd.Client.Revoke(context.TODO(), r.leaseID)
	return
}

func (r *Register) watch() (ok bool, err error) {
	regNode := &node2.Node{Name: r.node.Name}
	ok, err = etcd2.Get(regNode)
	if err != nil {
		return
	}
	return
}
