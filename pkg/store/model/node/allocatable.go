package node

import "time"

// Allocatable available resources for node
type Allocatable struct {
	NodeName   string    `json:"node_name"`
	Cpu        float64   `json:"cpu"`
	Disk       uint64    `json:"disk"`
	Memory     uint64    `json:"memory"`
	UpdateTime time.Time `json:"update_time"`
}

// AllocatableOption Allocatable's options
type AllocatableOption func(a *Allocatable)
