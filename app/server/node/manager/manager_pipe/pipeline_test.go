package manager_pipe

import (
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/store/model/scheduler"
	"testing"
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
