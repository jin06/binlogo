package election

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/jin06/binlogo/pkg/etcd_client"
	"github.com/jin06/binlogo/pkg/node/role"
	"github.com/jin06/binlogo/pkg/store/dao/dao_cluster"
	node2 "github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sync"
	"time"
)

type Election struct {
	node        *node2.Node
	election    *concurrency.Election
	client      *clientv3.Client
	campaignVal string
	ttl         int
	prefix      string
	role        role.Role
	lock        sync.Mutex
	roleMutex   sync.Mutex
	RoleCh      chan role.Role
}

func New(opts ...Option) (e *Election, err error) {
	e = &Election{
		ttl:    5,
		prefix: dao_cluster.ElectionPrefix(),
		role:   role.FOLLOWER,
		RoleCh: make(chan role.Role, 1000),
		lock:   sync.Mutex{},
	}
	for _, v := range opts {
		v(e)
	}
	err = e.init()
	return
}

func (e *Election) init() (err error) {
	e.campaignVal = viper.GetString("node.name")

	e.lock = sync.Mutex{}
	e.client = etcd_client.Default()
	return
}

func (e *Election) Run(ctx context.Context) {
	e.campaign(ctx)
}

func (e *Election) campaign(ctx context.Context) {
	go func() {
		defer func() {
			e.SetRole(role.FOLLOWER)
			e.Resign(ctx)
		}()
		for {
			e.SetRole(role.FOLLOWER)
			time.Sleep(time.Second)
			//sen, err := concurrency.NewSession(e.client, concurrency.WithTTL(e.ttl))
			//cli, err := etcd_client.New()
			//if err != nil {
			//	logrus.Error("campaign", err)
			//	continue
			//}

			sen, errSe := concurrency.NewSession(e.client, concurrency.WithTTL(e.ttl))
			if errSe != nil {
				logrus.Error("Election error")
				time.Sleep(time.Second * 5)
				continue
			}
			ele := concurrency.NewElection(
				sen,
				e.prefix,
			)
			e.election = ele

			logrus.Info("Run for election")
			ctx1, _ := context.WithCancel(ctx)
			errC := ele.Campaign(ctx1, e.campaignVal)
			if errC != nil {
				//_ = sen.Close()
				logrus.Error("campaign error", errC)
				time.Sleep(time.Second * 5)
				continue
			}

			logrus.Info("Win election")
			e.SetRole(role.LEADER)

			logrus.Info("Election Watch ")

			obsCh := e.election.Observe(ctx)
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
			if errResign := e.election.Resign(ctx); errResign != nil {
				logrus.Errorln("Election failed, start new election.", errResign)
			}
			logrus.Errorln("Lost leader... ")
			time.Sleep(time.Second)
		}
	}()
	return
}

func (e *Election) SetRole(r role.Role) {
	e.roleMutex.Lock()
	defer e.roleMutex.Unlock()
	e.role = r
	e.RoleCh <- r
}

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

func (e *Election) Role() role.Role {
	return e.role
}

func (e *Election) Resign(ctx context.Context) (err error) {
	if e.election == nil {
		return
	}
	err = e.election.Resign(ctx)
	return
}
