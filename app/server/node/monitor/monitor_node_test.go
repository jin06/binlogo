package monitor

import (
	"testing"

	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/pkg/store/dao/dao_node"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
	"github.com/jin06/binlogo/v2/pkg/util/random"
)

func TestCheckAllNode(t *testing.T) {
	configs.InitGoTest()
	nodeName := "gotest" + random.String()
	nModel := node.NewNode(nodeName)
	dao_node.CreateNode(nModel)
	defer dao_node.DeleteNode(nodeName)
	m, err := NewMonitor()
	if err != nil {
		t.Error(err)
	}
	err = m.checkAllNode()
	if err != nil {
		t.Error(err)
	}
}
