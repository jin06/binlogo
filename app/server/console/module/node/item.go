package node

import (
	"github.com/jin06/binlogo/pkg/node/role"
	"github.com/jin06/binlogo/pkg/store/model/node"
)

// Item for front display
type Item struct {
	Node        *node.Node        `json:"node"`
	Capacity    *node.Capacity    `json:"capacity"`
	Allocatable *node.Allocatable `json:"allocatable"`
	Info        *Info             `json:"info"`
	Status      *node.Status      `json:"status"`
}

// Info item info
type Info struct {
	Role role.Role `json:"role"`
}
