package node

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/pkg/store/model"
)

// Node is node
type Node struct {
	Name       string    `json:"name" redis:"name"`
	Role       Role      `json:"role" redis:"role"`
	IP         net.IP    `json:"ip" redis:"ip"`
	Version    string    `json:"version" redis:"version"`
	CreateTime time.Time `json:"create_time" redis:"create_time"`
}

// NewNode returns a new node
func NewNode(name string) *Node {
	return &Node{
		Name:       name,
		Version:    configs.Version,
		CreateTime: time.Now(),
	}
}

// Key generate etcd key
func (s *Node) Key() string {
	return fmt.Sprintf("node/%s", s.Name)
}

// Val generate etcd val
func (s *Node) Val() string {
	b, _ := json.Marshal(s)
	return string(b)
}

// Unmarshal generate from json data
func (s *Node) Unmarshal(data []byte) error {
	return json.Unmarshal(data, s)
}

// NodeOption configure node
type NodeOption func(s model.Values)

// WithNodeIP sets node's ip
func WithNodeIP(i net.IP) NodeOption {
	return func(s model.Values) {
		s["ip"] = i
	}
}

// WithNodeVersion sets node's version
func WithNodeVersion(v string) NodeOption {
	return func(s model.Values) {
		s["version"] = v
	}
}
