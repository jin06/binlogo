package node

type Node struct {
	Options *Options
}

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
