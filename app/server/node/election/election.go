package election

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jin06/binlogo/pkg/watcher"

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
	closeOnce   sync.Once
	// stopped
	// the channel "stopped" would be closed if election running process quit
	//
	stopped chan struct{}
}

// New returns a new Election
func New(opts ...Option) (e *Election) {
	e = &Election{
		ttl:    5,
		prefix: dao_cluster.ElectionPrefix(),
		role:   role.FOLLOWER,
		// RoleCh:  make(chan role.Role, 1000),
		eleLock: sync.Mutex{},
		stopped: make(chan struct{}),
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
func (e *Election) Run(ctx context.Context, ch chan role.Role) {
	defer e.close()
	e.RoleCh = ch
	e.campaign(ctx)
}

func (e *Election) campaign(ctx context.Context) {
	logrus.Info("Election cycle start")
	var err error
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln("election panic, ", r)
		}
		if err != nil {
			logrus.WithError(err).Errorln("election process quit unexpected")
		}
		logrus.Info("Election cycle end")
	}()
	client, err := etcdclient.New()
	if err != nil {
		return
	}
	if e.session, err = concurrency.NewSession(client, concurrency.WithTTL(e.ttl)); err != nil {
		return
	}
	e.election = concurrency.NewElection(e.session, e.prefix)

	logrus.Info("Run for election")
	errCh := make(chan error, 1)
	go func() {
		errCh <- e.election.Campaign(ctx, e.campaignVal)
	}()

	key := fmt.Sprintf("%s/%s/%x", dao_cluster.ElectionPrefix(), e.prefix, e.session.Lease())
	w, err := watcher.New(watcher.WithKey(key), watcher.WithHandler(watcher.WrapStrHandler()))
	if err != nil {
		return
	}
	ch, err := w.WatchEtcd(ctx)
	if err != nil {
		return
	}
	defer w.Close()
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
				s, er := dao_cluster.GetElection(fmt.Sprintf("%x", e.session.Lease()))
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
		case <-e.session.Done():
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

func (e *Election) close() (err error) {
	e.closeOnce.Do(func() {
		if e.session != nil {
			err = e.session.Close()
		}
		e.setRole(role.FOLLOWER)
		// close(e.RoleCh)
		close(e.stopped)
	})
	return
}

func (e *Election) Stopped() chan struct{} {
	return e.stopped
}
