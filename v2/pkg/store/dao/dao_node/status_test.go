package dao_node

import (
	"testing"

	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
	"github.com/jin06/binlogo/v2/pkg/util/random"
)

func TestStatus(t *testing.T) {
	configs.InitGoTest()
	nodeName := "go_test_node" + random.String()
	_, err := CreateOrUpdateStatus(nodeName,
		node.WithReady(true),
		node.WithCPUPressure(true),
		node.WithMemoryPressure(true),
		node.WithDiskPressure(true),
		node.WithNetworkUnavailable(true),
	)
	if err != nil {
		t.Error(err)
	}
	s, err := GetStatus(nodeName)
	if err != nil {
		t.Error(err)
	}
	if s.NodeName != nodeName {
		t.Log("name is different")
		t.Fail()
	}
	_, err = StatusMap()
	if err != nil {
		t.Error(err)
	}
	ok, err := DeleteStatus(nodeName)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Fail()
	}
}
