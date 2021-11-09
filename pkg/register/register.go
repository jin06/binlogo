package register

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/jin06/binlogo/pkg/etcd_client"
	"github.com/sirupsen/logrus"
	"time"
)

func New(opts ...Option) (r *Register, err error) {
	r = &Register{
		ttl: 5,
	}
	//r.client, err = etcd_client.New()
	r.client = etcd_client.Default()
	//if err != nil {
	//	return
	//}
	for _, v := range opts {
		v(r)
	}
	return
}

type Register struct {
	lease                  clientv3.Lease
	leaseID                clientv3.LeaseID
	ttl                    int64
	client                 *clientv3.Client
	registerKey            string
	registerData           interface{}
	registerCreateRevision int64
	ctx                    context.Context
}

func (r *Register) Run(ctx context.Context) {
	myCtx, c := context.WithCancel(ctx)
	r.ctx = myCtx
	go func() {
		defer func() {
			c()
		}()
		var err error
		for {
			err = r.reg()
			if err != nil {
				logrus.Errorln(err)
			}
		LOOP:
			for {
				select {
				case <-ctx.Done():
					{
						err = r.revoke()
						if err != nil {
							logrus.Errorln(err)
						}
						return
					}
				case <-time.Tick(time.Second):
					{
						err = r.keepOnce()
						if err == rpctypes.ErrLeaseNotFound {
							break LOOP
						}
						wOk, wErr := r.watch()
						if wErr != nil {
							logrus.Errorln(wErr)
						}
						if !wOk {
							break LOOP
						}
					}
				}
			}
			r.revoke()
		}
	}()
}

func (r *Register) reg() (err error) {
	//r.client, _ = etcd_client.New()
	r.lease = clientv3.NewLease(r.client)
	if rep, err2 := r.lease.Grant(context.TODO(), r.ttl); err2 != nil {
		return err2
	} else {
		r.leaseID = rep.ID
	}
	//err = dao_register.Reg(r.registerKey, r.registerData, r.leaseID)
	b, err := json.Marshal(r.registerData)
	if err != nil {
		return
	}
	res, err := r.client.Put(context.TODO(), r.registerKey, string(b), clientv3.WithLease(r.leaseID))
	if res != nil {
		r.registerCreateRevision = res.Header.Revision
	}
	return
}

func (r *Register) keepOnce() (err error) {
	_, err = r.lease.KeepAliveOnce(context.TODO(), r.leaseID)

	if err != nil {
		return
	}
	//logrus.Debugf("Node lease %v seconds.", resp.TTL)
	return
}

func (r *Register) revoke() (err error) {
	_, err = r.client.Revoke(context.TODO(), r.leaseID)
	return
}

func (r *Register) watch() (ok bool, err error) {
	var res *clientv3.GetResponse
	res, err = r.client.Get(context.TODO(), r.registerKey)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	if len(res.Kvs) == 0 {
		return
	}
	if res.Kvs[0].CreateRevision == r.registerCreateRevision {
		ok = true
	}
	return
}

func (r *Register) Context() context.Context {
	return r.ctx
}
