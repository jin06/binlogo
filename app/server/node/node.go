package node

import (
	"context"
	"github.com/jin06/binlogo/app/server/node/election"
	"github.com/jin06/binlogo/app/server/node/manager/manager_event"
	"github.com/jin06/binlogo/app/server/node/manager/manager_pipe"
	"github.com/jin06/binlogo/app/server/node/manager/manager_status"
	"github.com/jin06/binlogo/app/server/node/monitor"
	"github.com/jin06/binlogo/app/server/node/scheduler"
	"github.com/jin06/binlogo/pkg/node/role"
	"github.com/jin06/binlogo/pkg/register"
	store2 "github.com/jin06/binlogo/pkg/store"
	"github.com/jin06/binlogo/pkg/store/dao/dao_node"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

// Node represents a node instance
// Running pipeline, reporting status, etc.
// if it becomes the master node, it will run tasks such as scheduling pipeline, monitoring, event management, etc
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
	eventManager   *manager_event.Manager
}

type NodeMode byte

const (
	MODE_CLUSTER NodeMode = 1
	MODE_SINGLE  NodeMode = 2
)

// NodeMode todo
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

// New return a new node
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
		register.WithKey(dao_node.NodeRegisterPrefix()+"/"+n.Options.Node.Name),
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
	n.eventManager = manager_event.New()
	return
}

// Run start working
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

// Role returns current role
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
			case r := <-n.election.RoleCh:
				{
					n.leaderRun(ctx, r)
				}
			case <-time.Tick(time.Second * 5):
				{
					n.leaderRun(ctx, "")
				}
			}
		}
	}()
}

func (n *Node) leaderRun(ctx context.Context, r role.Role) {
	if r == "" {
		r = n.Role()
	}
	switch r {
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
			err = n.eventManager.Run(ctx)
			if err != nil {
				return
			}
		}
	default:
		{
			n.Scheduler.Stop()
			n.monitor.Stop(ctx)
			n.eventManager.Stop()
		}
	}
}
