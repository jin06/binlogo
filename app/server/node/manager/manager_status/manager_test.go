package manager_status

import (
	"context"
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"testing"
)

func TestRun(t *testing.T) {
	configs.DefaultEnv()
	configs.InitConfigs()
	m := NewManager(&node.Node{Name: "go_test_node"})
	err := m.Run(context.Background())
	if err != nil {
		t.Error(err)
	}
	ns := NewNodeStatus(m.Node.Name)
	if ns == nil {
		t.Fail()
	}
	err = ns.syncNodeStatus()
	if err != nil {
		t.Error(err)
	}

}
