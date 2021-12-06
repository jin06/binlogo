package election

import (
	"context"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/etcdclient"
	"github.com/jin06/binlogo/pkg/node/role"
	"github.com/jin06/binlogo/pkg/store/dao/dao_cluster"
	node2 "github.com/jin06/binlogo/pkg/store/model/node"
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
	//err = e.init()
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
		defer func() {
			e.SetRole(role.FOLLOWER)
			e.Resign(myCtx)
			cancel()
		}()
		d := time.Millisecond
		for {
			e.SetRole(role.FOLLOWER)
			//time.Sleep(time.Second)
			if e.election != nil {
				e.Resign(myCtx)
			}
			select {
			case <-ctx.Done():
				{
					return
				}
			case <-time.Tick(d):
				{
					d = time.Second * 5
				}
			}

			client, err := etcdclient.New()
			if err != nil {
				logrus.Errorln("Election new client error: ", err)
				return
			}
			sen, errSe := concurrency.NewSession(client, concurrency.WithTTL(e.ttl))
			if errSe != nil {
				logrus.Error("Election error")
				return
				//continue
			}
			ele := concurrency.NewElection(
				sen,
				e.prefix,
			)
			e.setElection(ele)

			logrus.Info("Run for election")
			errC := ele.Campaign(myCtx, e.campaignVal)
			if errC != nil {
				//_ = sen.Close()
				logrus.Error("campaign error", errC)
				return
				//continue
			}

			logrus.Info("Win election")
			e.SetRole(role.LEADER)

			logrus.Info("Election Watch ")

			obsCh := e.election.Observe(myCtx)
		LOOP:
			for {
				select {
				case <-ctx.Done():
					{
						return
					}
				case <-sen.Done():
					{
						break LOOP
					}
				case gr, ok := <-obsCh:
					{
						if !ok {
							break LOOP
						}
						if len(gr.Kvs) == 0 {
							break LOOP
						}
						if string(gr.Kvs[0].Value) != e.campaignVal {
							break LOOP
						}
					}
				}
			}
			if errResign := e.election.Resign(myCtx); errResign != nil {
				logrus.Errorln("Election failed, start new election.", errResign)
			}
			logrus.Errorln("Lost leader... ")
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

// getElection get etcd election client
func (e *Election) getElection() *concurrency.Election {
	e.eleLock.Lock()
	defer e.eleLock.Unlock()
	return e.election
}

// setElection sets etcd election client
func (e *Election) setElection(ele *concurrency.Election) {
	e.eleLock.Lock()
	defer e.eleLock.Unlock()
	e.election = ele
}

// Leader returns name of current leader node
func (e *Election) Leader() (name string, err error) {
	res, err := e.getElection().Leader(context.TODO())
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
