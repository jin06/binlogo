package node

import (
	"encoding/json"
	"github.com/jin06/binlogo/configs"
	"net"
	"time"
)

// Node is node
type Node struct {
	Name       string    `json:"name"`
	Role       Role      `json:"role"`
	IP         net.IP    `json:"ip"`
	Version    string    `json:"version"`
	CreateTime time.Time `json:"create_time"`
}

// NewNode returns a new node
func NewNode(name string) *Node {
	return &Node{
		Name:       name,
		Version:    configs.VERSITON,
		CreateTime: time.Now(),
	}
}

// Key generate etcd key
func (s *Node) Key() (key string) {
	key = "node/" + s.Name
	return
}

// Val generate etcd val
func (s *Node) Val() (val string) {
	b, _ := json.Marshal(s)
	val = string(b)
	return
}

// Unmarshal generate from json data
func (s *Node) Unmarshal(val []byte) (err error) {
	err = json.Unmarshal(val, s)
	return
}

// NodeOption configure node
type NodeOption func(s *Node)

// WithNodeIP sets node's ip
func WithNodeIP(i net.IP) NodeOption {
	return func(s *Node) {
		s.IP = i
	}
}

// WithNodeVersion sets node's version
func WithNodeVersion(v string) NodeOption {
	return func(s *Node) {
		s.Version = v
	}
}
