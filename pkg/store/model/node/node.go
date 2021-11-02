package node

import (
	"encoding/json"
	"net"
	"time"
)

type Node struct {
	Name       string    `json:"name"`
	Role       Role      `json:"role"`
	IP         net.IP    `json:"ip"`
	Version    string    `json:"version"`
	CreateTime time.Time `json:"create_time"`
}

func (s *Node) Key() (key string) {
	key = "node/" + s.Name
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
