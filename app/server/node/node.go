package node

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

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
)

// Node represents a node instance
// Running pipeline, reporting status, etc.
// if it becomes the master node, it will run tasks such as scheduling pipeline, monitoring, event management, etc
type Node struct {
	Mode     *NodeMode
	Options  *Options
	Name     string
	Register *register.Register
	//election       *election.Election
	electionManager *election.Manager
	Scheduler       *scheduler.Scheduler
	StatusManager   *manager_status.Manager
	monitor         *monitor.Monitor
	leaderRunMutex  sync.Mutex
	pipeManager     *manager_pipe.Manager
	eventManager    *manager_event.Manager
	// this channel watch all child process.
	// NOTED THIS: any process exist will call close(done) means 'one quit, all quit'
	// why? cos of all processes is essential, this simplified design is easy to impl, also well for debugging
	stopping chan struct{}
	stopOnce sync.Once
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
func New(opts ...Option) (node *Node) {
	options := &Options{}
	node = &Node{
		Options:  options,
		stopping: make(chan struct{}),
	}
	for _, v := range opts {
		v(options)
	}
	return
}

func (n *Node) init() {
	n.Register = register.New(
		register.WithTTL(5),
		register.WithKey(dao_node.NodeRegisterPrefix()+"/"+n.Options.Node.Name),
		register.WithData(n.Options.Node),
	)
	n.electionManager = election.NewManager(n.Options.Node)
	n.pipeManager = manager_pipe.New(n.Options.Node)
	n.StatusManager = manager_status.NewManager(n.Options.Node)
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
	stx, cancel := context.WithCancel(ctx)
	defer func() {
		if re := recover(); re != nil {
			err = errors.New(fmt.Sprintf("panic %v", re))
		}
		cancel()
	}()
	if err = n.refreshNode(); err != nil {
		return
	}

	n.init()

	go func() {
		n._mustRun(stx)
		n.stop()
	}()

	go func() {
		n._leaderRun(stx)
		n.stop()
	}()

	select {
	case <-ctx.Done():
		return
	case <-n.stopping:
		return
	}
}

// Role returns current role
func (n *Node) Role() role.Role {
	return n.electionManager.Role()
	//return n.election.Role()
}

// _leaderRun run when node is leader
func (n *Node) _leaderRun(ctx context.Context) {
	stx, cancel := context.WithCancel(ctx)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		case r := <-n.electionManager.RoleCh():
			{
				n.leaderRun(stx, r)
			}
		case <-time.Tick(time.Second * 3):
			{
				n.leaderRun(stx, n.Role())
			}
		}
	}
}

func (n *Node) leaderRun(ctx context.Context, r role.Role) {
	//if r == "" {
	//	r = n.Role()
	//}
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

// _mustRun run manager that every node must run
// such as election, node register, node status
func (n *Node) _mustRun(ctx context.Context) {
	stx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		if err := n.Register.Run(stx); err != nil {
			logrus.WithError(err).Errorln("node register exit")
		}
		n.stop()
	}()
	go func() {
		n.electionManager.Run(stx)
		n.stop()
	}()
	go func() {
		n.pipeManager.Run(stx)
		n.stop()
	}()
	go func() {
		n.StatusManager.Run(stx)
		n.stop()
	}()
	select {
	case <-ctx.Done():
	case <-n.stopping:
	}
}

func (n *Node) stop() {
	n.stopOnce.Do(
		func() {
			close(n.stopping)
		},
	)
}
