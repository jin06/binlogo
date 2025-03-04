package node

import (
	"github.com/jin06/binlogo/v2/internal/constant"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
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
	Role constant.Role `json:"role"`
}
