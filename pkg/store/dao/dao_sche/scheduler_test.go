package dao_sche

import (
	"testing"

	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/util/random"
)

func TestScheduler(t *testing.T) {
	configs.InitGoTest()
	pipeName := "go_test_pipeline" + random.String()
	nodeName := "go_test_node" + random.String()
	err := UpdatePipelineBindIfNotExist(pipeName, nodeName)
	if err != nil {
		t.Error(err)
	}
	_, err = UpdatePipelineBind(pipeName, nodeName)
	if err != nil {
		t.Error(err)
	}
	pb, err := GetPipelineBind()
	if err != nil {
		t.Error(err)
	}
	if pb == nil {
		t.Error("pipeline bind is nil")
	}

	ok, err := DeletePipelineBind(pipeName)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("delete not ok")
	}
}
