package event

import (
	"testing"

	"github.com/jin06/binlogo/configs"
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
	t.Log(NewInfoPipeline("go_test_pipeline", "test message"))
	t.Log(NewInfoCluster("test message"))
	t.Log(NewErrorPipeline("go_test_pipeline", "test message"))
	t.Log(NewErrorNode("test message"))
	t.Log(NewErrorCluster("test message"))
	t.Log(NewWarnPipeline("go_test_pipeline", "test message"))
	t.Log(NewWarnNode("test message"))
	t.Log(NewWarnCluster("test message"))

}
