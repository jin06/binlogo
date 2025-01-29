package node

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/jin06/binlogo/v2/app/server/node/election"
	"github.com/jin06/binlogo/v2/app/server/node/manager/manager_pipe"
	"github.com/jin06/binlogo/v2/app/server/node/manager/manager_status"
	"github.com/jin06/binlogo/v2/app/server/node/monitor"
	"github.com/jin06/binlogo/v2/app/server/node/scheduler"
	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/internal/constant"
	"github.com/jin06/binlogo/v2/pkg/store/dao"
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
	// eventManager    *manager_event.Manager
	// this channel watch all child process.
	// NOTED THIS: any process exist will call close(done) means 'one quit, all quit'
	// why? cos of all processes is essential, this simplified design is easy to impl, also well for debugging
	closing      chan struct{}
	closeOnce    sync.Once
	closed       chan struct{}
	completeOnce sync.Once
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
		Options: options,
		closing: make(chan struct{}),
		closed:  make(chan struct{}),
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

func (n *Node) refreshNode(ctx context.Context) error {
	if newest, err := dao.GetNode(ctx, n.Options.Node.Name); err != nil {
		return err
	} else {
		if newest == nil {
			return errors.New("unexpected, node is null")
		}
		n.Options.Node = newest
	}
	return nil
}

// Run start working
func (n *Node) Run(ctx context.Context) (err error) {
	defer n.CompleteClose()
	defer n.Close()
	defer func() {
		if re := recover(); re != nil {
			err = fmt.Errorf("Node run panic %v", re)
		}
	}()
	if err = n.refreshNode(ctx); err != nil {
		return
	}

	n.init()

	go func() {
		n.pollMustRun(ctx)
		n.Close()
	}()

	if configs.Default.Roles.Master {
		go func() {
			n.pollLeaderRun(ctx)
			n.Close()
		}()
	}

	select {
	case <-ctx.Done():
		return
	case <-n.closing:
		return
	}
}

// pollLeaderRun run when node is leader
func (n *Node) pollLeaderRun(ctx context.Context) {
	curRole := constant.FOLLOWER
	ticker := time.NewTicker(time.Second * 5)
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
				// curRole = n.electionManager.GetRole()
			}
		case <-n.closing:
			{
				return
			}
		}
		n.leaderRun(ctx, curRole)
	}
}

func (n *Node) leaderRun(ctx context.Context, r constant.Role) {
	n.leaderRunMutex.Lock()
	defer n.leaderRunMutex.Unlock()
	//logrus.Debug("node run leader ", r)
	if r == constant.FOLLOWER {
		if n.Scheduler != nil {
			n.Scheduler.Close()
			n.Scheduler = nil
		}
		if n.monitor != nil {
			n.monitor.Close()
			n.monitor = nil
		}
		// if n.eventManager != nil {
		// 	n.eventManager.Stop()
		// 	n.eventManager = nil
		// }
	}
	if r == constant.LEADER {
		if n.Scheduler == nil || n.Scheduler.IsClosed() {
			n.Scheduler = scheduler.New()
			go n.Scheduler.Run(ctx)
		}
		if n.monitor == nil || n.monitor.IsClosed() {
			n.monitor = monitor.NewMonitor()
			go n.monitor.Run(ctx)
		}
		// if n.eventManager == nil || n.eventManager.Exit {
		// n.eventManager = manager_event.New()
		// go n.eventManager.Run(ctx)
		// }
	}
}

// pollMustRun run manager that every node must run
// such as election, node register, node status
func (n *Node) pollMustRun(ctx context.Context) {
	// defer n.CompleteClose()
	// defer n.Close()
	stx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		if err := n.Register.Run(stx); err != nil {
			slog.Error("node register exit: " + err.Error())
		}
		n.Close()
	}()
	if configs.Default.Roles.Master {
		go func() {
			n.electionManager.Run(stx)
			n.Close()
		}()
	}
	go func() {
		n.pipeManager.Run(stx)
		n.Close()
	}()
	go func() {
		n.StatusManager.Run(stx)
		n.Close()
	}()

	select {
	case <-ctx.Done():
	case <-n.closing:
	}
}

func (n *Node) Closed() chan struct{} {
	return n.closed
}

func (n *Node) CompleteClose() {
	n.completeOnce.Do(func() {
		<-n.monitor.Closed()
		<-n.pipeManager.Closed()
		<-n.electionManager.Closed()
		<-n.Register.Closed()
		close(n.closed)
	})
}

func (n *Node) Close() {
	n.closeOnce.Do(func() {
		close(n.closing)
	})
}
