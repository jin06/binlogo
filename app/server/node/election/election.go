package election

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/etcdclient"
	"github.com/jin06/binlogo/pkg/node/role"
	"github.com/jin06/binlogo/pkg/store/dao/dao_cluster"
	node2 "github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/watcher/str"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

// Election for leader node election
type Election struct {
	node        *node2.Node
	election    *concurrency.Election
	campaignVal string
	ttl         int
	prefix      string
	role        role.Role
	eleLock     sync.Mutex
	roleMutex   sync.Mutex
	RoleCh      chan role.Role
	ctx         context.Context
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
	e.campaign(ctx)
}

func (e *Election) campaign(ctx context.Context) {
	myCtx, cancel := context.WithCancel(ctx)
	e.ctx = myCtx
	go func() {
		logrus.Info("Election cycle start")
		var err error
		defer func() {
			e.SetRole(role.FOLLOWER)
			if err != nil {
				logrus.Error("Election failed: ", err)
			}
			errR := e.Resign(ctx)
			if errR != nil {
				logrus.Error(errR)
			}
			cancel()
			logrus.Info("Election cycle end")
		}()
		e.SetRole(role.FOLLOWER)
		client, err := etcdclient.New()
		if err != nil {
			return
		}
		sen, err := concurrency.NewSession(client, concurrency.WithTTL(e.ttl))
		if err != nil {
			return
		}
		ele := concurrency.NewElection(
			sen,
			e.prefix,
		)
		e.election = ele

		logrus.Info("Run for election")
		errCh := make(chan error, 1)
		go func() {
			errCh <- ele.Campaign(ctx, e.campaignVal)
		}()
		ch, err := str.Watch(ctx, dao_cluster.ElectionPrefix(), fmt.Sprintf("%s/%x", e.prefix, sen.Lease()))
		if err != nil {
			return
		}
		for {
			var b bool
			select {
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
				close(ch)
				break
			}
		}

		logrus.Info("Win election")
		e.SetRole(role.LEADER)

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
	}()
	return
}

// SetRole sets node current role
func (e *Election) SetRole(r role.Role) {
	e.roleMutex.Lock()
	defer e.roleMutex.Unlock()
	e.role = r
	e.RoleCh <- r
}

// Leader returns name of current leader node
func (e *Election) Leader() (name string, err error) {
	res, err := e.election.Leader(context.TODO())
	if err != nil {
		return
	}
	if len(res.Kvs) == 0 {
		return
	}
	name = string(res.Kvs[0].Value)
	return
}

// Role returns role current node
func (e *Election) Role() role.Role {
	return e.role
}

// Resign resign election
func (e *Election) Resign(ctx context.Context) (err error) {
	if e.election == nil {
		return
	}
	err = e.election.Resign(ctx)
	if err == nil {
		e.SetRole(role.FOLLOWER)
	}
	return
}

// Context returns Election context
func (e *Election) Context() context.Context {
	return e.ctx
}
