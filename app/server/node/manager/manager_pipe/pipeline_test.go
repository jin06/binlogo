package manager_pipe

import (
	"context"
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/store/dao/dao_sche"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/store/model/scheduler"
	"testing"
	"time"
)

func TestScanPipelines(t *testing.T) {
	pb := &scheduler.PipelineBind{
		Bindings: map[string]string{
			"apple":  "test",
			"banana": "test2",
			"orange": "test",
		},
	}
	m := New(&node.Node{Name: "test"})
	m.mapping["tomato"] = true

	err := m.scanPipelines(pb)
	if err != nil {
		t.Error(err)
	}
	if m.mapping["apple"] == false {
		t.Fail()
	}
	if m.mapping["banana"] == true {
		t.Fail()
	}
	if m.mapping["orange"] == false {
		t.Fail()
	}
	if m.mapping["tomato"] == true {
		t.Fail()
	}
}

func TestRun(t *testing.T) {
	configs.DefaultEnv()
	configs.InitConfigs()
	_, err := dao_sche.UpdatePipelineBind("go_test_pipeline", "go_test_node")
	if err != nil {
		t.Error(err)
	}
	m := New(&node.Node{
		Name: "go_test_node",
	})
	m.Run(context.Background())
	time.Sleep(time.Millisecond*200)
}
