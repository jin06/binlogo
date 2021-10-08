package model

import (
	"encoding/json"
	"net"
)

type NodeRole string

const (
	ROLE_MASTER NodeRole = "master"
	ROLE_SLAVE  NodeRole = "slave"
)

type Node struct {
	Name string   `json:"name"`
	Role NodeRole `json:"role"`
	Ip   net.IP  `json:"ip"`
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
func (s *Node) Unmarshal(val []byte) (err error) {
	err = json.Unmarshal(val, s)
	return
}

