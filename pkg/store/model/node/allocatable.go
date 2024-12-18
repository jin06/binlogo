package node

import (
	"encoding/json"
	"fmt"
	"time"
)

// Allocatable available resources for node
type Allocatable struct {
	NodeName   string    `json:"node_name" redis:"node_name"`
	Cpu        float64   `json:"cpu" redis:"cpu"`
	Disk       uint64    `json:"disk" redis:"disk"`
	Memory     uint64    `json:"memory" redis:"memory"`
	UpdateTime time.Time `json:"update_time" redis:"update_time"`
}

func (a Allocatable) MarshalBinary() (data []byte, err error) {
	// encoding.BinaryMarshaler
	return json.Marshal(a)
}

func (a *Allocatable) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, a)
}

func (a *Allocatable) UnmarshalText(data []byte) error {
	return json.Unmarshal(data, a)
}

func (a *Allocatable) Key() string {
	return fmt.Sprintf("node/allocatable/%s", a.NodeName)
}

func (a *Allocatable) Val() string {
	b, _ := json.Marshal(a)
	return string(b)
}

func (a *Allocatable) Unmarshal(data []byte) error {
	return json.Unmarshal(data, a)
}

// AllocatableOption Allocatable's options
type AllocatableOption func(a *Allocatable)

// NewAllocatable returns a new *Allocatable
func NewAllocatable(nodeName string) *Allocatable {
	a := &Allocatable{
		NodeName:   nodeName,
		Cpu:        0,
		Disk:       0,
		Memory:     0,
		UpdateTime: time.Time{},
	}
	return a
}
