package register

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/etcd_client"
	"github.com/jin06/binlogo/pkg/store/dao/dao_cluster"
	node2 "github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/store/model/register"
	"github.com/sirupsen/logrus"
	"time"
)

func New(opts ...Option) (r *Register, err error) {
	r = &Register{
		ttl: 5,
	}
	r.client, err = etcd_client.New()
	logrus.Debug("Node register lease duration ", r.ttl)
	for _, v := range opts {
		v(r)
	}
	logrus.Debug("Node register lease duration ", r.ttl)
	return
}

type Register struct {
	lease   clientv3.Lease
	leaseID clientv3.LeaseID
	node    *node2.Node
	ttl     int64
	client  *clientv3.Client
}

func (r *Register) Run(ctx context.Context) (err error) {
	go func() {
		for {
			err = r.reg()
			if err != nil {
				logrus.Errorln(err)
			}
			for {
				select {
				case <-time.Tick(time.Second):
					{
						err = r.keepOnce()
						if err != nil {
							break
						}
					}
				case <-ctx.Done():
					{
						er := r.revoke()
						if er != nil {
							logrus.Errorln(er)
						}
						return
					}
				}
			}

		}
	}()

	return
}

func (r *Register) reg() (err error) {
	r.lease = clientv3.NewLease(r.client)
	if rep, err2 := r.lease.Grant(context.TODO(), r.ttl); err2 != nil {
		return err2
	} else {
		r.leaseID = rep.ID
	}
	err = dao_cluster.RegRNode(&register.RegisterNode{Name: r.node.Name}, r.leaseID)
	return
}

func (r *Register) keepOnce() (err error) {
	resp, err := r.lease.KeepAliveOnce(context.TODO(), r.leaseID)
	if err != nil {
		//logrus.Error(err)
		return err
	}
	logrus.Debugf("Node lease %v seconds.", resp.TTL)
	return
}
func (r *Register) revoke() (err error) {
	_, err = r.client.Revoke(context.TODO(), r.leaseID)
	return
}

func (r *Register) watch() (ok bool, err error) {
	var regNode *register.RegisterNode
	regNode, err = dao_cluster.GetRNode(r.node.Name)
	if err != nil {
		blog.Error("get rnode error")
		return
	}
	if regNode != nil {
		ok = true
	}
	return
}
