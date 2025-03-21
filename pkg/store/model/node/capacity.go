package node

import (
	"encoding/json"
	"fmt"
	"time"
)

type CapacityMap map[string]*Capacity

func (c *CapacityMap) Key() string {
	return "/capacity_map"
}

func (c *CapacityMap) Val() string {
	b, _ := json.Marshal(c)
	return string(b)
}

func (c *CapacityMap) Unmarshal(data []byte) error {
	return json.Unmarshal(data, c)
}

// Capacity node capacity
type Capacity struct {
	NodeName    string      `json:"node_name" redis:"node_name"`
	Cpu         float64     `json:"cpu" redis:"cpu"`
	Disk        uint64      `json:"disk" redis:"disk"`
	Memory      uint64      `json:"memory" redis:"memory"` //byte
	CpuCors     int32       `json:"cpu_cores" redis:"cpu_cores"`
	CpuUsage    uint8       `json:"cpu_usage" redis:"cpu_usage"`
	MemoryUsage uint8       `json:"memory_usage" redis:"memory_usage"`
	UpdateTime  time.Time   `json:"update_time" redis:"update_time"`
	Allocatable Allocatable `json:"allocatable" redis:"allocatable"`
}

func (c *Capacity) Key() string {
	return fmt.Sprintf("node_capacity/%s", c.NodeName)
}

func (c *Capacity) Val() string {
	b, _ := json.Marshal(c)
	return string(b)
}

func (c *Capacity) Unmarshal(data []byte) error {
	return json.Unmarshal(data, c)
}

// New returns a new *Capacity
func NewCapacity(nodeName string) *Capacity {
	c := &Capacity{}
	return c
}

// CapacityOption sets Capacity
type CapacityOption func(c *Capacity)
