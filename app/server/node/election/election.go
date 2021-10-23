package election

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/node/role"
	node2 "github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.etcd.io/etcd/clientv3/concurrency"
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
		ttl: 1,
		prefix: fmt.Sprintf(
			"/%v/%v/election",
			configs.APP,
			viper.GetString("cluster.name"),
		),
		role:   role.FOLLOWER,
		RoleCh: make(chan role.Role, 1000),
	}
	for _, v := range opts {
		v(e)
	}
	err = e.Init()
	return
}

func (e *Election) Init() (err error) {
	e.client, err = clientv3.New(clientv3.Config{
		Endpoints:   viper.GetStringSlice("store.etcd.endpoints"),
		DialTimeout: 2 * time.Second,
	})
	//ipVal, err := ip.LocalIp()
	//if err != nil {
	//	return
	//}
	//e.campaignVal = ipVal.String()
	e.campaignVal = viper.GetString("node.name")

	e.lock = sync.Mutex{}
	return
}

func (e *Election) Run(ctx context.Context) {
	e.campaign(ctx)
	go func() {
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			}
		}
	}()
}

func (e *Election) campaign(ctx context.Context) {
	go func() {
		for {
			e.SetRole(role.FOLLOWER)
			sen, err := concurrency.NewSession(e.client, concurrency.WithTTL(e.ttl))
			if err != nil {
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
			err = ele.Campaign(ctx1, e.campaignVal)
			if err != nil {
				_ = sen.Close()
				time.Sleep(time.Second * 5)
				continue
			}

			logrus.Info("Win election")
			e.SetRole(role.LEADER)

			logrus.Info("Election Watch ")
			ctx2, _ := context.WithCancel(ctx)
			//obsCh := ele.Observe(ctx2)
			for {
				time.Sleep(time.Second)
				res, err := ele.Leader(ctx2)
				if err == concurrency.ErrElectionNoLeader {
					logrus.Info("No leader ...")
					break
				}
				//fmt.Printf("%v", res.Kvs)
				for _, v := range res.Kvs {
					if e.campaignVal == string(v.Value) {
						logrus.Info("I'm leader")
						//roleCh <- role.LEADER
						break
					} else {
						logrus.Info("I'm not leader")
						e.SetRole(role.FOLLOWER)
					}
				}

			}
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

func (e *Election) Role() role.Role {
	return e.role
}

func (e *Election) Resign(ctx context.Context) (err error) {
	err = e.election.Resign(ctx)
	return
}
