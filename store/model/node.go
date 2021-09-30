package model

import "encoding/json"

type NodeRole string

const (
	ROLE_MASTER NodeRole = "master"
	ROLE_SLAVE  NodeRole = "slave"
)

type Node struct {
	Name string   `json:"name"`
	Role NodeRole `json:"role"`
}

func (s *Node) Key() (key string) {
	key = "nodes/" + s.Name
	return
}

func (s *Node) Val() (val string) {
	b, _ := json.Marshal(s)
	val = string(b)
	return
}
