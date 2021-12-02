package node

import (
	"context"
	"errors"
	"github.com/jin06/binlogo/app/server/node/election"
	"github.com/jin06/binlogo/app/server/node/manager/manager_event"
	"github.com/jin06/binlogo/app/server/node/manager/manager_pipe"
	"github.com/jin06/binlogo/app/server/node/manager/manager_status"
	"github.com/jin06/binlogo/app/server/node/monitor"
	"github.com/jin06/binlogo/app/server/node/scheduler"
	"github.com/jin06/binlogo/pkg/node/role"
	"github.com/jin06/binlogo/pkg/register"
	"github.com/jin06/binlogo/pkg/store/dao/dao_node"
	"github.com/jin06/binlogo/pkg/store/model/node"
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
	myCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	var newest *node.Node
	newest, err = dao_node.GetNode(n.Options.Node.Name)
	if err != nil {
		return
	}
	if newest == nil {
		err = errors.New("unexpected, node is null")
		return
	}
	n.Options.Node = newest
	n.Register.Run(myCtx)
	err = n.StatusManager.Run(myCtx)
	n.election = election.New(
		election.OptionNode(n.Options.Node),
		election.OptionTTL(5),
	)
	n.election.Run(myCtx)
	n.pipeManager.Run(myCtx)
	n._leaderRun(myCtx)
	select {
	case <-ctx.Done():
		{
			return
		}
	case <-n.Register.Context().Done():
		{
			return
		}
	case <-n.election.Context().Done():
		{
			return
		}
	}
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
