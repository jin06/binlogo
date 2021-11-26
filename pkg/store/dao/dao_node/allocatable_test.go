package dao_node

import (
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/util/random"
	"testing"
	"time"
)

func TestUpdateAllocatable(t *testing.T) {
	configs.InitGoTest()
	nodeName := "go_test_node" + random.String()
	al := node.Allocatable{
		NodeName:   nodeName,
		Cpu:        0,
		Disk:       0,
		Memory:     0,
		UpdateTime: time.Time{},
	}
	err := UpdateAllocatable(&al)
	if err != nil {
		t.Error(err)
	}
	_, err = DeleteAllocatable(nodeName)
	if err != nil {
		t.Error(err)
	}
}
