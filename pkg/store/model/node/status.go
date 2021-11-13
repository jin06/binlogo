package node

type Status struct {
	Ready              bool `json:"ready"`
	NetworkUnavailable bool `json:"network_unavailable"`
	MemoryPressure     bool `json:"memory_pressure"`
	DiskPressure       bool `json:"disk_pressure"`
	CPUPressure        bool `json:"cpu_pressure"`
}

type StatusOption func(s *Status)

func WithReady(b bool) StatusOption {
	return func(s *Status) {
		s.Ready = b
	}
}

func WithNetworkUnavailable( b bool)  StatusOption{
	return func(s *Status) {
		s.NetworkUnavailable = b
	}
}

func WithMemoryPressure( b bool) StatusOption  {
	return func(s *Status) {
		s.MemoryPressure = b
	}
}

func WithDiskPressure(b bool) StatusOption  {
	return func(s *Status) {
		s.DiskPressure = b
	}
}

func WithCPUPressure(b bool ) StatusOption  {
	return func(s *Status) {
		s.CPUPressure = b
	}
}
