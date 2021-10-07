package node

import (
	register2 "github.com/jin06/binlogo/app/server/node/register"
	"github.com/jin06/binlogo/store"
	"time"
)

type Node struct {
	Role     *NodeRole
	Mode     *NodeMode
	Options  *Options
	Name     string
	Register *register2.Register
}

type NodeRole byte

const (
	ROLE_MASTER NodeRole = 1
	ROLE_SLAVE  NodeRole = 2
)

func (n NodeRole) String() string {
	switch n {
	case ROLE_MASTER:
		{
			return "master"
		}
	case ROLE_SLAVE:
		{
			return "slave"
		}

	}
	return ""
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
	return
}

func (n *Node) Run() (err error) {
	ok, err := store.Get(n.Options.Node)
	if err != nil {
		return
	}
	if ok {
		//todo
		panic("123")
	}
	err = n.Register.Run()
	if err != nil {
		return
	}

	select {}
	return
}
