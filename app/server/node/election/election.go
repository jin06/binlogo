package election

import (
	"context"
	"fmt"
	"github.com/jin06/binlogo/pkg/watcher"
	"sync"
	"time"

	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/etcdclient"
	"github.com/jin06/binlogo/pkg/node/role"
	"github.com/jin06/binlogo/pkg/store/dao/dao_cluster"
	node2 "github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3/concurrency"
)

// Election for leader node election
type Election struct {
	node        *node2.Node
	session     *concurrency.Session
	election    *concurrency.Election
	campaignVal string
	ttl         int
	prefix      string
	role        role.Role
	eleLock     sync.Mutex
	roleMutex   sync.Mutex
	RoleCh      chan role.Role
	ctx         context.Context
	closeOnce   sync.Once
	cancel      context.CancelFunc
}

// New returns a new Election
func New(opts ...Option) (e *Election) {
	e = &Election{
		ttl:     5,
		prefix:  dao_cluster.ElectionPrefix(),
		role:    role.FOLLOWER,
		RoleCh:  make(chan role.Role, 1000),
		eleLock: sync.Mutex{},
	}
	for _, v := range opts {
		v(e)
	}
	if e.node == nil {
		e.campaignVal = configs.NodeName
	} else {
		e.campaignVal = e.node.Name
	}
	return
}

// Run start election process
func (e *Election) Run(ctx context.Context) {
	e.ctx, e.cancel = context.WithCancel(ctx)
	go func() {
		defer e.close()
		e.campaign(e.ctx)
	}()
}

func (e *Election) campaign(ctx context.Context) {
	logrus.Info("Election cycle start")
	var err error
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln("election panic, ", r)
			panic(r)
		}
		//e.setRole(role.FOLLOWER)
		//if err != nil {
		//	logrus.Error("Election failed: ", err)
		//}
		//errR := e.resign(e.ctx)
		//if errR != nil {
		//	logrus.Error(errR)
		//}
		logrus.Info("Election cycle end")
	}()
	client, err := etcdclient.New()
	if err != nil {
		return
	}
	sen, err := concurrency.NewSession(client, concurrency.WithTTL(e.ttl))
	if err != nil {
		return
	}
	e.session = sen
	e.election = concurrency.NewElection(
		sen,
		e.prefix,
	)

	logrus.Info("Run for election")
	errCh := make(chan error, 1)
	go func() {
		errCh <- e.election.Campaign(e.ctx, e.campaignVal)
	}()

	key := fmt.Sprintf("%s/%s/%x", dao_cluster.ElectionPrefix(), e.prefix, sen.Lease())
	w, err := watcher.New(watcher.WithKey(key), watcher.WithHandler(watcher.WrapStrHandler()))
	if err != nil {
		return
	}
	ch, err := w.WatchEtcd(e.ctx)
	if err != nil {
		return
	}
	defer func() {
		w.Close()
	}()
	for {
		var b bool
		select {
		case <-ctx.Done():
			{
				return
			}
		case err = <-errCh:
			{
				if err != nil {
					return
				} else {
					b = true
					break
				}
			}
		case ev, ok := <-ch:
			{
				if !ok {
					return
				}
				if ev.Event.Type == mvccpb.DELETE {
					return
				}
			}
		case <-time.Tick(time.Minute):
			{
				s, er := dao_cluster.GetElection(fmt.Sprintf("%x", sen.Lease()))
				if er != nil {
					return
				}
				if s == nil {
					return
				}
			}
		}
		if b {
			break
		}
	}

	logrus.Info("Win election")
	e.setRole(role.LEADER)

	logrus.Info("Election Watch ")

	obsCh := e.election.Observe(ctx)
	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		case <-sen.Done():
			{
				return
			}
		case gr, ok := <-obsCh:
			{
				if !ok {
					return
				}
				if len(gr.Kvs) == 0 {
					return
				}
				if string(gr.Kvs[0].Value) != e.campaignVal {
					return
				}
			}
		}
	}
}

// SetRole sets node current role
func (e *Election) setRole(r role.Role) {
	e.roleMutex.Lock()
	defer e.roleMutex.Unlock()
	e.role = r
	e.RoleCh <- r
}

// Role returns role current node
func (e *Election) getRole() role.Role {
	return e.role
}

// resign resign election
// start new election
func (e *Election) resign(ctx context.Context) (err error) {
	if e.election == nil {
		return
	}
	err = e.election.Resign(ctx)
	//if err == nil {
	//	e.setRole(role.FOLLOWER)
	//}
	return
}

// Context returns Election context
func (e *Election) context() context.Context {
	return e.ctx
}

func (e *Election) close() (err error) {
	e.closeOnce.Do(func() {
		err = e.session.Close()
		close(e.RoleCh)
		e.cancel()
	})
	return
}
