package dao_node

import (
	"github.com/jin06/binlogo/configs"
	node2 "github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/util/random"
	"testing"
	"time"
)

func TestNode(t *testing.T) {
	configs.InitGoTest()
	nodeName := "go_test_node" + random.String()
	nModel := &node2.Node{
		Name:       nodeName,
		Role:       node2.Role{},
		IP:         nil,
		Version:    "",
		CreateTime: time.Time{},
	}
	err := CreateNodeIfNotExist(nModel)
	if err != nil {
		t.Error(err)
	}
	err = CreateNode(nModel)
	if err != nil {
		t.Error(err)
	}
	_, err = UpdateNode(nodeName, node2.WithNodeVersion(configs.VERSITON))
	if err != nil {
		t.Error(err)
	}
	nMode2, err := GetNode(nodeName)
	if err != nil {
		t.Error(err)
	}
	if nMode2.Version != configs.VERSITON {
		t.Fail()
	}
	if nMode2.Name != nModel.Name {
		t.Fail()
	}
	list, err := AllNodes()
	if err != nil {
		t.Error()
	}
	if len(list) < 0 {
		t.Fail()
	}
	_, err = ALLRegisterNodes()
	if err != nil {
		t.Error(err)
	}
	_, err = AllNodesMap()
	if err != nil {
		t.Error(err)
	}
	_, err = AllRegisterNodesMap()
	if err != nil {
		t.Error(err)
	}
	_, err = AllWorkNodesMap()
	if err != nil {
		t.Error(err)
	}
	ok, err := DeleteNode(nodeName)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Fail()
	}
}
