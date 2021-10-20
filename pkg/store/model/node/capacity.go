package node

type Capacity struct {
	Cpu    float64 `json:"cpu"`
	Disk   uint64 `json:"disk"`
	Memory uint64 `json:"memory"` //byte
}
