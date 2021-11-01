package node

type Status struct {
	Ready              bool `json:"ready"`
	NetworkUnavailable bool `json:"network_unavailable"`
	MemoryPressure     bool `json:"memory_pressure"`
	DiskPressure       bool `json:"disk_pressure"`
	CPUPressure        bool `json:"cpu_pressure"`
}
