package node

import (
	"context"
	"github.com/jin06/binlogo/app/server/node/election"
	"github.com/jin06/binlogo/app/server/node/manager/manager_pipe"
	"github.com/jin06/binlogo/app/server/node/manager/manager_status"
	"github.com/jin06/binlogo/app/server/node/monitor"
	"github.com/jin06/binlogo/app/server/node/scheduler"
	"github.com/jin06/binlogo/pkg/node/role"
	"github.com/jin06/binlogo/pkg/register"
	store2 "github.com/jin06/binlogo/pkg/store"
	"github.com/jin06/binlogo/pkg/store/dao/dao_cluster"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Node struct {
	Mode           *NodeMode
	Options        *Options
	Name           string
	Register       *register.Register
	election       *election.Election
	Scheduler      *scheduler.Scheduler
	StatusManager  *manager_status.Manager
	monitor        *monitor.Monitor
	leaderRunMutex sync.Mutex
	pipeManager    *manager_pipe.Manager
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
		Options: options,
	}
	for _, v := range opts {
		v(options)
	}
	err = node.init()
	return
}

func (n *Node) init() (err error) {
	n.Register, err = register.New(
		register.WithTTL(5),
		register.WithKey(dao_cluster.RegisterPrefix() + "/" + n.Options.Node.Name),
		register.WithData(n.Options.Node),
	)
	if err != nil {
		return
	}
	n.Scheduler = scheduler.New()

	logrus.Debug("---->", n.Options.Node)
	n.StatusManager = manager_status.NewManager(n.Options.Node)
	n.monitor, err = monitor.NewMonitor()
	if err != nil {
		return
	}
	n.pipeManager = manager_pipe.New(n.Options.Node)
	return
}

func (n *Node) Run(ctx context.Context) (err error) {
	go func() {
		var ok bool
		if ok, err = store2.Get(n.Options.Node); err != nil {
			return
		}
		if ok {
			//todo
			//panic("exist node")
		}
		ctxReg, cancelReg := context.WithCancel(ctx)
		defer cancelReg()
		n.Register.Run(ctxReg)

		ctxStatus, cancelStatus := context.WithCancel(ctx)
		defer cancelStatus()
		err = n.StatusManager.Run(ctxStatus)

		ctxElection, cancelElection := context.WithCancel(ctx)
		defer cancelElection()
		n.election, err = election.New(
			election.OptionNode(n.Options.Node),
			election.OptionTTL(5),
		)
		if err != nil {
			return
		}
		n.election.Run(ctxElection)

		n.pipeManager.Run(ctx)

		ctxLeader, cancelLeader := context.WithCancel(ctx)
		n._leaderRun(ctxLeader)
		defer cancelLeader()

		select {
		case <-ctx.Done():
			{
			panic(ctx.Err())
				return
			}
		}
	}()
	return
}

func (n *Node) Role() role.Role {
	return n.election.Role()
}

func (n *Node) _leaderRun(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			case <-n.election.RoleCh:
				{
					n.leaderRun(ctx)
				}
			case <-time.Tick(time.Second * 2):
				{
					n.leaderRun(ctx)
				}
			}
		}
	}()
}

func (n *Node) leaderRun(ctx context.Context) {
	switch n.Role() {
	case role.LEADER:
		{
			var err error
			err = n.Scheduler.Run(ctx)
			if err != nil {
				return
			}
			err = n.monitor.Run(ctx)
			if err != nil {
				return
			}
		}
	default:
		{
			n.Scheduler.Stop(ctx)
			n.monitor.Stop(ctx)
		}
	}
}
