package node

import "github.com/jin06/binlogo/pkg/store/model/node"

type Item struct {
	Node *node.Node `json:"node"`
	Info *Info      `json:"info"`
}

type Info struct {
}
