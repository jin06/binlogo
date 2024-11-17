package node

import "encoding/json"

// Status node status
type Status struct {
	NodeName           string `json:"node_name" redis:"node_name"`
	Ready              bool   `json:"ready" redis:"ready"`
	NetworkUnavailable bool   `json:"network_unavailable" redis:"network_unavailable"`
	MemoryPressure     bool   `json:"memory_pressure" redis:"memory_pressure"`
	DiskPressure       bool   `json:"disk_pressure" redis:"disk_pressure"`
	CPUPressure        bool   `json:"cpu_pressure" redis:"cpu_pressure"`
}

func (s *Status) Key() string {
	return "node/status/" + s.NodeName
}

func (s *Status) Val() string {
	b, _ := json.Marshal(s)
	return string(b)
}

func (s *Status) Unmarshal(data []byte) error {
	return json.Unmarshal(data, s)
}

// New returns a status model
func New(nodeName string) (s *Status) {
	s = &Status{
		NodeName:           nodeName,
		Ready:              true,
		NetworkUnavailable: false,
		MemoryPressure:     false,
		DiskPressure:       false,
		CPUPressure:        false,
	}
	return
}

// StatusOption is function configure Status
type StatusOption func(s *Status)

// WithReady sets status ready
func WithReady(b bool) StatusOption {
	return func(s *Status) {
		s.Ready = b
	}
}

// WithNetworkUnavailable sets status network
func WithNetworkUnavailable(b bool) StatusOption {
	return func(s *Status) {
		s.NetworkUnavailable = b
	}
}

// WithMemoryPressure sets status memory
func WithMemoryPressure(b bool) StatusOption {
	return func(s *Status) {
		s.MemoryPressure = b
	}
}

// WithDiskPressure sets status disk
func WithDiskPressure(b bool) StatusOption {
	return func(s *Status) {
		s.DiskPressure = b
	}
}

// WithCPUPressure sets status cpu
func WithCPUPressure(b bool) StatusOption {
	return func(s *Status) {
		s.CPUPressure = b
	}
}
