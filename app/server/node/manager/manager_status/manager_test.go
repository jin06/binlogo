package manager_status

import (
	"context"
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/store/dao/dao_node"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/util/random"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	configs.DefaultEnv()
	configs.InitConfigs()
	nodeName := "go_test_node" + random.String()
	nModel := node.NewNode(nodeName)
	dao_node.CreateNode(nModel)
	defer dao_node.DeleteNode(nodeName)
	m := NewManager(&node.Node{Name: nodeName})
	ctx, cancel := context.WithCancel(context.Background())
	err := m.Run(ctx)
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
	err = m.syncStatus()
	if err != nil {
		t.Error(err)
	}
	err = m.syncIP()
	if err != nil {
		t.Error(err)
	}
	cancel()
	time.Sleep(time.Millisecond * 20)
}
