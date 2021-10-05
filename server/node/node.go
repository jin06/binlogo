package node

type Node struct {
	Role    *NodeRole
	Options *Options
}

type NodeRole byte

var (
	RoleMaster = 1
	RoleSlave  = 2
)

func New(opts ...Option) (node *Node, err error) {
	options := &Options{}
	node = &Node{
		Options: options,
	}
	for _, v := range opts {
		v(options)
	}
	return
}
