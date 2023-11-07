package register

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	"sync"
	"time"

	"github.com/jin06/binlogo/pkg/etcdclient"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// New returns a new *Register
func New(opts ...Option) (r *Register) {
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
	cancel                 context.CancelFunc
	closeOnce              sync.Once
}

func (r *Register) initClient() (err error) {
	//if r.client != nil {
	//	er := r.client.Close()
	//	if er != nil {
	//		logrus.Errorln(er)
	//	}
	//}
	r.client, err = etcdclient.New()
	return
}

// Run start register
func (r *Register) Run(ctx context.Context) {
	logrus.Info("register run ")
	r.ctx, r.cancel = context.WithCancel(ctx)
	go func() {
		defer func() {
			if re := recover(); re != nil {
				logrus.Errorln("register panic, ", re)
			}
			r.close()
		}()
		var err error
		if r.client, err = etcdclient.New(); err != nil {
			logrus.Errorln(err)
			return
		}
		err = r.reg()
		if err != nil {
			logrus.Errorln(err)
			return
		}
		defer func() {
			errR := r.revoke()
			if errR != nil {
				logrus.Errorln(errR)
			}
			logrus.Errorln("Register end: ", r.registerKey)
		}()

		keepTicker := time.NewTicker(time.Second)
		defer keepTicker.Stop()
		watchTicker := time.NewTicker(time.Second * 5)
		defer watchTicker.Stop()
		for {
			select {
			case <-r.ctx.Done():
				return
			case <-ctx.Done():
				{
					return
				}
			case <-watchTicker.C:
				{
					wOk, wErr := r.watch()
					if wErr != nil {
						logrus.Errorln(wErr)
						return
					}
					if !wOk {
						return
					}
				}
			case <-keepTicker.C:
				{
					err = r.keepOnce()
					if err != nil {
						fmt.Println(err)
					}
					if err == rpctypes.ErrLeaseNotFound {
						return
					}
				}
			}
		}
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
	txn := r.client.Txn(r.ctx).
		If(clientv3.Compare(clientv3.CreateRevision(r.registerKey), "=", int64(0))).
		Then(clientv3.OpPut(r.registerKey, string(b), clientv3.WithLease(r.leaseID)))
	res, err := txn.Commit()
	if err != nil {
		return
	}
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
		//logrus.Errorln(err)
		return
	}
	if len(res.Kvs) == 0 {
		return
	}
	if res.Kvs[0].CreateRevision == r.registerCreateRevision {
		ok = true
	}
	//logrus.Debug(res.Kvs[0].CreateRevision)
	logrus.Debug(r.registerCreateRevision)
	return
}

// Context returns register's context
func (r *Register) Context() context.Context {
	return r.ctx
}

func (r *Register) close() error {
	r.closeOnce.Do(func() {
		r.cancel()
		r.client.Close()
	})
	return nil
}
