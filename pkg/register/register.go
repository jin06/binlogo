package register

import (
	"context"
	"encoding/json"
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
		ttl:     5,
		stopped: make(chan struct{}),
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
	stopped                chan struct{}
	closeOnce              sync.Once
}

func (r *Register) init() (err error) {
	logrus.Info("init register")
	if r.client, err = etcdclient.New(); err != nil {
		return
	}
	return
}

// Run start register
func (r *Register) Run(ctx context.Context) (err error) {
	logrus.WithField("key", r.registerKey).Info("register run")
	defer func() {
		if re := recover(); re != nil {
			logrus.Errorln("register panic, ", re)
		}
		if err != nil {
			logrus.WithError(err).Errorln("register process quit unexpectedly")
		}
		logrus.WithField("key", r.registerKey).Info("register stopped")
		r.close()
	}()
	if err = r.init(); err != nil {
		return
	}
	stx, cancel := context.WithCancel(ctx)
	defer cancel()
	if err = r.reg(stx); err != nil {
		logrus.Errorln(err)
		return
	}
	defer func() {
		errR := r.revoke(stx)
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
		case <-ctx.Done():
			{
				return
			}
		case <-watchTicker.C:
			{
				wOk, wErr := r.watch(stx)
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
				err = r.keepOnce(stx)
				if err == rpctypes.ErrLeaseNotFound {
					return
				}
			}
		}
	}
}

func (r *Register) reg(ctx context.Context) (err error) {
	//r.client, _ = etcdclient.New()
	r.lease = clientv3.NewLease(r.client)
	rep, err := r.lease.Grant(ctx, r.ttl)
	if err != nil {
		return
	}
	r.leaseID = rep.ID
	//err = dao_register.Reg(r.registerKey, r.registerData, r.leaseID)
	b, err := json.Marshal(r.registerData)
	if err != nil {
		return
	}
	txn := r.client.Txn(ctx).
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

func (r *Register) keepOnce(ctx context.Context) (err error) {
	_, err = r.lease.KeepAliveOnce(ctx, r.leaseID)

	if err != nil {
		return
	}
	//logrus.Debugf("Node lease %v seconds.", resp.TTL)
	return
}

func (r *Register) revoke(ctx context.Context) (err error) {
	_, err = r.client.Revoke(ctx, r.leaseID)
	return
}

func (r *Register) watch(ctx context.Context) (ok bool, err error) {
	var res *clientv3.GetResponse
	res, err = r.client.Get(ctx, r.registerKey)
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

func (r *Register) close() error {
	r.closeOnce.Do(func() {
		if r.client != nil {
			r.client.Close()
		}
		close(r.stopped)
	})
	return nil
}

func (r *Register) Stopped() chan struct{} {
	return r.stopped
}
