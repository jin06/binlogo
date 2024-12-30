package manager_status

import (
	"context"
	"errors"
	"time"

	"github.com/jin06/binlogo/v2/pkg/store/dao"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/sirupsen/logrus"
)

// NodeStatus reporting work
type NodeStatus struct {
	NodeName    string
	Node        *node.Node       `json:"node" redis:"node"`
	Cap         node.Capacity    `json:"cap" redis:"cap"`
	Allocatable node.Allocatable `json:"allocatable" redis:"allocatable"`
	Conditions  *node.Condition  `json:"conditions" redis:"conditions"`
}

// NewNodeStatus returns a new NodeStatus
func NewNodeStatus(nodeName string) *NodeStatus {
	ns := &NodeStatus{}
	ns.NodeName = nodeName
	return ns
}

func (ns *NodeStatus) syncNodeStatus() (err error) {
	// logrus.Infoln("Sync node status ")
	err = ns.setStatus()
	if err != nil {
		return
	}
	err = ns.syncCap()
	if err != nil {
		logrus.Errorf("sync cap error: %s", err)
		return
	}
	err = ns.syncAllocatable()
	if err != nil {
		logrus.Errorf("sync allocatable error: %s", err)
		return
	}
	return
}

func (ns *NodeStatus) setStatus() (err error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return
	}
	ts, err := cpu.Times(false)
	if err != nil {
		return
	}
	if len(ts) <= 0 {
		err = errors.New("cpu check error")
		return
	}

	c := ts[0]

	capacity := node.Capacity{
		NodeName:   ns.NodeName,
		Memory:     v.Total,
		Cpu:        c.User + c.System + c.Idle,
		UpdateTime: time.Now(),
	}
	if cpuInfo, err1 := cpu.Info(); err1 == nil {
		var cores int32
		for _, info := range cpuInfo {
			cores = cores + info.Cores
		}
		capacity.CpuCors = cores
	} else {
		logrus.Error(err1)
	}

	al := node.Allocatable{
		NodeName:   ns.NodeName,
		Memory:     v.Available,
		Cpu:        c.Idle,
		UpdateTime: time.Now(),
	}

	//capacity.CpuUsage, _ = decimal.NewFromFloat(al.Cpu/capacity.Cpu).Round(2).Float64()
	capacity.CpuUsage = uint8((capacity.Cpu - al.Cpu) * 100 / capacity.Cpu)

	capacity.MemoryUsage = uint8((capacity.Memory - al.Memory) * 100 / capacity.Memory)
	capacity.Allocatable = al

	ns.Cap = capacity
	ns.Allocatable = al
	return
}

func (ns *NodeStatus) syncCap() (err error) {
	_, err = dao.UpdateCapacity(context.Background(), &ns.Cap)
	return
}

func (ns *NodeStatus) syncAllocatable() (err error) {
	_, err = dao.UpdateAllocatable(context.Background(), &ns.Allocatable)
	return
}
