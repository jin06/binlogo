package node

import (
	"context"
	"github.com/jin06/binlogo/app/server/node/election"
	register2 "github.com/jin06/binlogo/app/server/node/register"
	"github.com/jin06/binlogo/app/server/node/scheduler"
	"github.com/jin06/binlogo/pkg/node/role"
	store2 "github.com/jin06/binlogo/pkg/store"
	"github.com/jin06/binlogo/pkg/watcher/wathcer_pipeline"
	"github.com/sirupsen/logrus"
	"time"
)

type Node struct {
	Mode      *NodeMode
	Options   *Options
	Name      string
	Register  *register2.Register
	election  *election.Election
	Scheduler *scheduler.Scheduler
}

type NodeMode byte

const (
	MODE_CLUSTER NodeMode = 1
	MODE_SINGLE  NodeMode = 2
)

func (n NodeMode) String() string {
	switch n {
	case MODE_CLUSTER:
		{
			return "cluster"
		}
	case MODE_SINGLE:
		{
			return "single"
		}
	}
	return ""
}

func New(opts ...Option) (node *Node, err error) {
	options := &Options{}
	node = &Node{
		Options:  options,
		Register: &register2.Register{},
	}
	for _, v := range opts {
		v(options)
	}
	err = node.Init()
	return
}

func (n *Node) Init() (err error) {
	n.Register = register2.New(
		register2.OptionNode(n.Options.Node),
		register2.OptionLeaseDuration(2*time.Second),
	)
	n.Scheduler = scheduler.New()
	return
}

func (n *Node) Run(ctx context.Context) (err error) {
	ok, err := store2.Get(n.Options.Node)

	if err != nil {
		return
	}
	if ok {
		//todo
		//panic("exist node")
	}
	ctx1, cancel1 := context.WithCancel(ctx)
	defer cancel1()
	err = n.Register.Run(ctx1)
	if err != nil {
		return
	}

	ctx2, cancel2 := context.WithCancel(ctx)
	defer cancel2()
	n.election, err = election.New(
		election.OptionNode(n.Options.Node),
		election.OptionTTL(5),
	)
	if err != nil {
		return
	}
	n.election.Run(ctx2)
	defer cancel2()

	ctx3, cancel3 := context.WithCancel(ctx)
	n.scheduler(ctx3)
	defer cancel3()

	logrus.Debug("watch")
	w, err := pipeline.New("/binlogo/cluster1/pipeline/test")
	if err != nil {
		logrus.Error(err)
		return
	}
	w.Watch(ctx)

	select {}
	return
}

func (n *Node) Role() role.Role {
	return n.election.Role()
}

func (n *Node) scheduler(ctx context.Context) {
	go func() {
		for range time.Tick(time.Second * 5) {
			logrus.Debug("scheduler ...")
			if n.Role() == role.LEADER {
				n.Scheduler.Run(ctx)
			} else {
				n.Scheduler.Stop(ctx)
			}
		}
	}()
}
