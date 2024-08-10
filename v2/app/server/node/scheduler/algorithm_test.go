package scheduler

import (
	"testing"
	"time"

	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
	"github.com/jin06/binlogo/v2/pkg/store/model/scheduler"
)

func TestCal(t *testing.T) {
	nodeName1 := "go_test_node_1"
	node1 := node.Node{
		Name:       nodeName1,
		Version:    configs.VERSITON,
		CreateTime: time.Now(),
	}
	nodeName2 := "go_test_node_2"
	node2 := node.Node{
		Name:       nodeName2,
		Version:    configs.VERSITON,
		CreateTime: time.Now(),
	}
	cap1 := node.Capacity{
		NodeName:    nodeName1,
		Cpu:         7369562.59,
		Disk:        0,
		Memory:      8589934592,
		CpuCors:     4,
		CpuUsage:    19,
		MemoryUsage: 72,
		UpdateTime:  time.Time{},
		Allocatable: &node.Allocatable{
			NodeName:   nodeName1,
			Cpu:        5970494.52,
			Disk:       0,
			Memory:     1950904320,
			UpdateTime: time.Time{},
		},
	}
	cap2 := node.Capacity{
		NodeName:    nodeName2,
		Cpu:         7369562.59,
		Disk:        0,
		Memory:      8589934592,
		CpuCors:     4,
		CpuUsage:    19,
		MemoryUsage: 72,
		UpdateTime:  time.Time{},
		Allocatable: &node.Allocatable{
			NodeName:   nodeName2,
			Cpu:        5970494.52,
			Disk:       0,
			Memory:     1950904320,
			UpdateTime: time.Time{},
		},
	}
	schePipeName := "go_test_pipeline_now"
	pipeName1 := "go_test_pipeline_1"
	alg := newAlgorithm(
		withAlgPipe(&pipeline.Pipeline{Name: schePipeName}),
		withAlgAllNodes(map[string]*node.Node{
			nodeName1: &node1,
			nodeName2: &node2,
		}),
		withAlgPipeBind(&scheduler.PipelineBind{Bindings: map[string]string{
			pipeName1: nodeName1,
		}}),
		withAlgCapMap(map[string]*node.Capacity{
			nodeName1: &cap1,
			nodeName2: &cap2,
		}),
	)
	err := alg.cal()
	if err != nil {
		t.Fail()
	}
	if alg.bestNode.Name == nodeName1 {
		t.Log(alg.bestNode)
		t.Fail()
	}
}
