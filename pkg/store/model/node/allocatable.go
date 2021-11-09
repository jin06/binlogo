package node

import "time"

type Allocatable struct {
	Cpu    float64 `json:"cpu"`
	Disk   uint64  `json:"disk"`
	Memory uint64  `json:"memory"`
	UpdateTime time.Time `json:"update_time"`
}
