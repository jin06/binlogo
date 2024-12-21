package node

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/jin06/binlogo/v2/app/server/node/election"
	"github.com/jin06/binlogo/v2/app/server/node/manager/manager_event"
	"github.com/jin06/binlogo/v2/app/server/node/manager/manager_pipe"
	"github.com/jin06/binlogo/v2/app/server/node/manager/manager_status"
	"github.com/jin06/binlogo/v2/app/server/node/monitor"
	"github.com/jin06/binlogo/v2/app/server/node/scheduler"
	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/internal/constant"
	"github.com/jin06/binlogo/v2/pkg/store/dao"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
)

// Node represents a node instance
// Running pipeline, reporting status, etc.
// if it becomes the master node, it will run tasks such as scheduling pipeline, monitoring, event management, etc
type Node struct {
	Mode    *NodeMode
	Options *Options
	Name    string
	// Register *register.Register
	Register *Register
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
	// main goroutine exited: all related goroutines exited
}

type NodeMode byte

const (
	MODE_CLUSTER NodeMode = 1
	MODE_SINGLE  NodeMode = 2
)

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
	n.Register = NewRegister(n.Options.Node)
	n.electionManager = election.NewManager(n.Options.Node)
	n.pipeManager = manager_pipe.New(n.Options.Node)
	n.StatusManager = manager_status.NewManager(n.Options.Node)
}

func (n *Node) refreshNode() (err error) {
	var newest *node.Node
	newest, err = dao.GetNode(context.Background(), n.Options.Node.Name)
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
			err = fmt.Errorf("panic %v", re)
		}
		cancel()
		n.stop()
	}()
	if err = n.refreshNode(); err != nil {
		return
	}

	n.init()

	go func() {
		n.pollMustRun(stx)
	}()

	if configs.Default.Roles.Master {
		go func() {
			n.pollLeaderRun(stx)
		}()
	}

	select {
	case <-ctx.Done():
		return
	case <-n.stopping:
		return
	}
}

// pollLeaderRun run when node is leader
func (n *Node) pollLeaderRun(ctx context.Context) {
	stx, cancel := context.WithCancel(ctx)
	defer func() {
		defer cancel()
		n.stop()
	}()
	curRole := constant.FOLLOWER
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-n.electionManager.Closed():
			return
		case r, ok := <-n.electionManager.RoleCh():
			{
				if !ok {
					return
				}
				curRole = r
			}
		case <-ticker.C:
			{

			}
		case <-n.stopping:
			{
				return
			}
		}
		n.leaderRun(stx, curRole)
	}
}

func (n *Node) leaderRun(ctx context.Context, r constant.Role) {
	n.leaderRunMutex.Lock()
	defer n.leaderRunMutex.Unlock()
	//logrus.Debug("node run leader ", r)
	if r == constant.FOLLOWER {
		if n.Scheduler != nil {
			n.Scheduler.Stop()
			n.Scheduler = nil
		}
		if n.monitor != nil {
			n.monitor.Stop()
			n.monitor = nil
		}
		if n.eventManager != nil {
			n.eventManager.Stop()
			n.eventManager = nil
		}
	}
	if r == constant.LEADER {
		if n.Scheduler == nil || n.Scheduler.Exit {
			n.Scheduler = scheduler.New()
			go n.Scheduler.Run(ctx)
		}
		if n.monitor == nil || n.monitor.Exit {
			n.monitor = monitor.NewMonitor()
			go n.monitor.Run(ctx)
		}
		if n.eventManager == nil || n.eventManager.Exit {
			n.eventManager = manager_event.New()
			go n.eventManager.Run(ctx)
		}
	}
}

// pollMustRun run manager that every node must run
// such as election, node register, node status
func (n *Node) pollMustRun(ctx context.Context) {
	stx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		if err := n.Register.Run(stx); err != nil {
			slog.Error("node register exit: " + err.Error())
		}
		n.stop()
	}()
	if configs.Default.Roles.Master {
		go func() {
			n.electionManager.Run(stx)
			n.stop()
		}()
	}
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
	n.stopOnce.Do(func() {
		close(n.stopping)
	})
}
