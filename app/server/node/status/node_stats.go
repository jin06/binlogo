package status

import (
	"errors"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/store/dao"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type NodeStatus struct {
	NodeName    string
	Node        *node.Node
	Cap         *node.Capacity
	Allocatable *node.Allocatable
	Conditions  *node.Condition
}

func NewNodeStatus(nodeName string) *NodeStatus {
	ns := &NodeStatus{}
	ns.NodeName = nodeName
	return ns
}

func (ns *NodeStatus) syncNodeStatus() (err error) {
	blog.Debug("Sync node status ")
	err = ns.setStatus()
	if err != nil {
		return
	}
	err = ns.syncCap()
	if err != nil {
		return
	}
	err = ns.syncAllocatable()
	if err != nil {
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

	capacity := &node.Capacity{
		Memory: v.Total,
		Cpu:    c.User + c.System + c.Idle,
	}
	if cpuInfo, err1 := cpu.Info(); err1 == nil {
		var cores int32
		for _, info := range cpuInfo {
			cores = cores + info.Cores
		}
		capacity.CpuCors = cores
	} else {
		blog.Error(err1)
	}

	al := &node.Allocatable{
		Memory: v.Available,
		Cpu:    c.Idle,
	}

	//capacity.CpuUsage, _ = decimal.NewFromFloat(al.Cpu/capacity.Cpu).Round(2).Float64()
	capacity.CpuUsage = uint8((capacity.Cpu - al.Cpu) * 100 / capacity.Cpu)

	capacity.MemoryUsage = uint8((capacity.Memory - al.Memory) * 100 / capacity.Memory)

	ns.Cap = capacity
	ns.Allocatable = al
	return
}

func (ns *NodeStatus) syncCap() (err error) {
	//key := etcd.Prefix() + "/nodes/" + ns.NodeName + "/capacity"
	//err = dao.UpdateCapacity(ns.Cap, dao.WithKey(key))
	err = dao.UpdateCapacity(ns.Cap, dao.WithNodeName(ns.NodeName))
	return
}

func (ns *NodeStatus) syncAllocatable() (err error) {
	//key := etcd.Prefix() + "/nodes/" + ns.NodeName + "/allocatable"
	//err = dao.UpdateAllocatable(ns.Allocatable, dao.WithKey(key))
	err = dao.UpdateAllocatable(ns.Allocatable, dao.WithNodeName(ns.NodeName))
	return
}
