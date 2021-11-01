package node

import (
	"encoding/json"
	"github.com/jin06/binlogo/pkg/node/role"
	"net"
	"time"
)

type NodeStatus string

const (
	STATUS_ON  = "on"
	STATUS_OFF = "off"
)

type Node struct {
	Name       string     `json:"name"`
	Role       role.Role  `json:"role"`
	IP         net.IP     `json:"ip"`
	Status     NodeStatus `json:"status"`
	Version    string     `json:"version"`
	CreateTime time.Time  `json:"create_time"`
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
