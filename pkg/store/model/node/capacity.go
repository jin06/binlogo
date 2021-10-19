package node

type Capacity struct {
	Cpu    uint   `json:"cpu"`
	Disk   uint64 `json:"disk"`
	Memory uint   `json:"memory"`
}
