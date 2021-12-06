package node

import (
	"context"
	"errors"
	"github.com/jin06/binlogo/app/server/node/election"
	"github.com/jin06/binlogo/app/server/node/manager"
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
	logrus.Debug("---->", n.Options.Node)
	return
}

func (n *Node) refreshNode() (err error) {
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
	return
}

// Run start working
func (n *Node) Run(ctx context.Context) (err error) {
	myCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	err = n.refreshNode()
	if err != nil {
		return
	}
	n.Register = register.New(
		register.WithTTL(5),
		register.WithKey(dao_node.NodeRegisterPrefix()+"/"+n.Options.Node.Name),
		register.WithData(n.Options.Node),
	)
	n.Register.Run(myCtx)
	n.election = election.New(
		election.OptionNode(n.Options.Node),
		election.OptionTTL(5),
	)
	n.election.Run(myCtx)
	n.pipeManager = manager_pipe.New(n.Options.Node)
	n.pipeManager.Run(myCtx)
	n.StatusManager = manager_status.NewManager(n.Options.Node)
	err = n.StatusManager.Run(myCtx)
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
			if n.Scheduler == nil || n.Scheduler.Status() == scheduler.SCHEDULER_STOP {
				n.Scheduler = scheduler.New()
				err = n.Scheduler.Run(ctx)
				if err != nil {
					return
				}
			}
			if n.monitor == nil || n.monitor.Status() == monitor.STATUS_STOP {
				n.monitor, err = monitor.NewMonitor()
				if err != nil {
					return
				}
				err = n.monitor.Run(ctx)
				if err != nil {
					return
				}
			}
			if n.eventManager == nil || n.eventManager.Status == manager.STOP {
				n.eventManager = manager_event.New()
				err = n.eventManager.Run(ctx)
				if err != nil {
					return
				}

			}
		}
	default:
		{
			if n.Scheduler != nil {
				n.Scheduler.Stop()
			}
			if n.monitor != nil {
				n.monitor.Stop(ctx)
			}
			if n.eventManager != nil {
				n.eventManager.Stop()
			}
		}
	}
}
