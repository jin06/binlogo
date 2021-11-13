package node

import "time"

type Capacity struct {
	NodeName    string    `json:"node_name"`
	Cpu         float64   `json:"cpu"`
	Disk        uint64    `json:"disk"`
	Memory      uint64    `json:"memory"` //byte
	CpuCors     int32     `json:"cpu_cors"`
	CpuUsage    uint8     `json:"cpu_usage"`
	MemoryUsage uint8     `json:"memory_usage"`
	UpdateTime  time.Time `json:"update_time"`
	Allocatable *Allocatable
}

type CapacityOption func(c *Capacity)
