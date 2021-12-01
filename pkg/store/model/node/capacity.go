package node

import "time"

// Capacity node capacity
type Capacity struct {
	NodeName    string       `json:"node_name"`
	Cpu         float64      `json:"cpu"`
	Disk        uint64       `json:"disk"`
	Memory      uint64       `json:"memory"` //byte
	CpuCors     int32        `json:"cpu_cores"`
	CpuUsage    uint8        `json:"cpu_usage"`
	MemoryUsage uint8        `json:"memory_usage"`
	UpdateTime  time.Time    `json:"update_time"`
	Allocatable *Allocatable `json:"allocatable"`
}

// New returns a new *Capacity
func NewCapacity(nodeName string) *Capacity {
	c := &Capacity{
		NodeName:    nodeName,
		Cpu:         0,
		Disk:        0,
		Memory:      0,
		CpuCors:     0,
		CpuUsage:    0,
		MemoryUsage: 0,
		UpdateTime:  time.Time{},
		Allocatable: NewAllocatable(nodeName),
	}
	return c
}

// CapacityOption sets Capacity
type CapacityOption func(c *Capacity)
