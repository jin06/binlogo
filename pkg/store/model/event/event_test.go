package event

import (
	"github.com/jin06/binlogo/configs"
	"testing"
)

func TestNewInfoNode(t *testing.T) {
	nodeName := "go_test_node"
	configs.NodeName = nodeName
	e := NewInfoNode("just test")
	if e.Type != INFO {
		t.Fail()
	}
	if e.ResourceType != NODE {
		t.Fail()
	}
	if e.NodeName != nodeName {
		t.Fail()
	}
}
