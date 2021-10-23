package node

import (
	"context"
	"github.com/jin06/binlogo/app/server/node/election"
	register2 "github.com/jin06/binlogo/app/server/node/register"
	"github.com/jin06/binlogo/app/server/node/scheduler"
	"github.com/jin06/binlogo/app/server/node/status"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/node/role"
	store2 "github.com/jin06/binlogo/pkg/store"
	"sync"
	"time"
)

type Node struct {
	Mode           *NodeMode
	Options        *Options
	Name           string
	Register       *register2.Register
	election       *election.Election
	Scheduler      *scheduler.Scheduler
	StatusManager  *status.Manager
	leaderRunMutex sync.Mutex
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
	n.StatusManager = &status.Manager{
		NodeName: n.Options.Node.Name,
	}
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
		if err = n.Register.Run(ctxReg); err != nil {
			return
		}

		ctxElection, cancelElection := context.WithCancel(ctx)
		defer cancelElection()
		if n.election, err = election.New(
			election.OptionNode(n.Options.Node),
			election.OptionTTL(5),
		); err != nil {
			return
		}
		n.election.Run(ctxElection)

		ctxLeader, cancelLeader := context.WithCancel(ctx)
		n._leaderRun(ctxLeader)
		defer cancelLeader()

		ctxStatus, cancelStatus := context.WithCancel(ctx)
		defer cancelStatus()
		err = n.StatusManager.Run(ctxStatus)
		select {
		case <-ctx.Done():
			{
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
	if n.Role() == role.LEADER {
		if err := n.Scheduler.Run(ctx); err != nil {
			blog.Error("Scheduler run failed : ", err)
			err = n.election.Resign(ctx)
			if err != nil {
				blog.Error("Election resign failed: ")
			}
		}
	} else {
		n.Scheduler.Stop(ctx)
	}
}
