package dao_node

import (
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/util/random"
	"testing"
	"time"
)

func TestCapacity(t *testing.T) {
	configs.InitGoTest()
	nodeName := "go_test_node" + random.String()
	c := &node.Capacity{
		NodeName:    nodeName,
		Cpu:         0,
		Disk:        0,
		Memory:      0,
		CpuCors:     0,
		CpuUsage:    0,
		MemoryUsage: 0,
		UpdateTime:  time.Time{},
		Allocatable: nil,
	}
	err := UpdateCapacity(c)
	if err != nil {
		t.Error(err)
	}
	cm, err := CapacityMap()
	if err != nil {
		t.Error()
	}
	if len(cm) < 0 {
		t.Fail()
	}
	_, err = DeleteCapacity(nodeName)
	if err != nil {
		t.Error(err)
	}
}
