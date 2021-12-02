package register

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/jin06/binlogo/pkg/etcdclient"
	"github.com/sirupsen/logrus"
	"time"
)

// New returns a new *Register
func New(opts ...Option) (r *Register, err error) {
	r = &Register{
		ttl: 5,
	}
	//r.client = etcdclient.Default()
	for _, v := range opts {
		v(r)
	}
	return
}

// Register Encapsulation of etcd register
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

func (r *Register) initClient() (err error) {
	if r.client != nil {
		er := r.client.Close()
		if er != nil {
			logrus.Errorln(er)
		}
	}
	r.client, err = etcdclient.New()
	return
}

// Run start register
func (r *Register) Run(ctx context.Context) {
	myCtx, cancel := context.WithCancel(ctx)
	r.ctx = myCtx
	go func() {
		defer func() {
			cancel()
		}()
		var err error
		err = r.initClient()
		if err != nil {
			logrus.Errorln(err)
			return
		}
		err = r.reg()
		if err != nil {
			logrus.Errorln(err)
			return
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
		errR := r.revoke()
		if errR != nil {
			logrus.Errorln(errR)
		}
		logrus.Errorln("Register end: ", r.registerKey)
		//time.Sleep(time.Second)
	}()
}

func (r *Register) reg() (err error) {
	//r.client, _ = etcdclient.New()
	r.lease = clientv3.NewLease(r.client)
	rep, err := r.lease.Grant(r.ctx, r.ttl)
	if err != nil {
		return
	}
	r.leaseID = rep.ID
	//err = dao_register.Reg(r.registerKey, r.registerData, r.leaseID)
	b, err := json.Marshal(r.registerData)
	if err != nil {
		return
	}
	res, err := r.client.Put(r.ctx, r.registerKey, string(b), clientv3.WithLease(r.leaseID))
	if res != nil {
		r.registerCreateRevision = res.Header.Revision
	}
	return
}

func (r *Register) keepOnce() (err error) {
	_, err = r.lease.KeepAliveOnce(r.ctx, r.leaseID)

	if err != nil {
		return
	}
	//logrus.Debugf("Node lease %v seconds.", resp.TTL)
	return
}

func (r *Register) revoke() (err error) {
	_, err = r.client.Revoke(r.ctx, r.leaseID)
	return
}

func (r *Register) watch() (ok bool, err error) {
	var res *clientv3.GetResponse
	res, err = r.client.Get(r.ctx, r.registerKey)
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

// Context returns register's context
func (r *Register) Context() context.Context {
	return r.ctx
}
