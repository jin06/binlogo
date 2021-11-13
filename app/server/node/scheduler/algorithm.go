package scheduler

import (
	"errors"
	"github.com/jin06/binlogo/pkg/store/dao/dao_node"
	"github.com/jin06/binlogo/pkg/store/dao/dao_sche"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/jin06/binlogo/pkg/store/model/scheduler"
)

type algorithm struct {
	pipeline       *pipeline.Pipeline
	allNodes       map[string]*node.Node
	potentialNodes map[string]*node.Node
	nodesScores    map[string]int
	bestNode       *node.Node
	pb             *scheduler.PipelineBind
	capacityMap    map[string]*node.Capacity
}

func newAlgorithm(p *pipeline.Pipeline) *algorithm {
	a := &algorithm{}
	a.pipeline = p
	a.allNodes = map[string]*node.Node{}
	a.potentialNodes = map[string]*node.Node{}
	a.nodesScores = map[string]int{}
	a.bestNode = &node.Node{}
	return a
}

func (a *algorithm) cal() (err error) {
	a.allNodes, err = dao_node.AllWorkNodesMap()
	if err != nil {
		return
	}
	err = a.calPotentialNodes()
	if err != nil {
		return err
	}
	err = a.calScore()
	if err != nil {
		return
	}
	err = a.calBestNode()
	if err != nil {
		return
	}
	return
}

func (a *algorithm) calPotentialNodes() (err error) {
	a.potentialNodes = a.allNodes
	return
}

func (a *algorithm) calBestNode() (err error) {
	if len(a.potentialNodes) == 0 {
		err = errors.New("no potential node")
		return
	}
	score := 0

	for k, _ := range a.potentialNodes {
		if a.nodesScores[k] > score {
			a.bestNode = a.potentialNodes[k]
			score = a.nodesScores[k]
		}
	}
	return
}

func (a *algorithm) calScore() (err error) {
	a.pb, err = dao_sche.GetPipelineBind()
	if err != nil {
		return
	}
	a.capacityMap, err = dao_node.CapacityMap()
	if err != nil {
		return
	}
	a.nodesScores = map[string]int{}
	for _, v := range a.potentialNodes {
		a.nodesScores[v.Name] = 0
	}
	scores, err := a._scoreNumOfPipelines()
	if err != nil {
		return
	}
	a.mergeScores(scores)
	scores, err = a._scoreResources()
	if err != nil {
		return
	}
	a.mergeScores(scores)
	return
}

func (a *algorithm) mergeScores(scores map[string]int) {
	if scores == nil {
		return
	}
	for pName, score := range scores {
		if _, ok := a.nodesScores[pName]; ok {
			a.nodesScores[pName] += score
		}
	}
	return
}

// Node's score depend on numbers of running pipeline in the node.
func (a *algorithm) _scoreNumOfPipelines() (scores map[string]int, err error) {
	weight := 5
	scores = map[string]int{}
	pipeNums := map[string]int{}
	totalNode := 0
	totalPipe := 0
	for _, v := range a.potentialNodes {
		pipeNums[v.Name] = 0
		totalNode++
	}
	for _, v := range a.pb.Bindings {
		totalPipe++
		if _, ok := pipeNums[v]; ok {
			pipeNums[v]++
		}
	}
	var numberPerNode float64
	numberPerNode = float64(totalPipe) / float64(totalNode)
	for k, v := range pipeNums {
		score := 0
		var per float64
		if totalPipe > 0 {
			per = float64(v) / numberPerNode
			if per <= 0 {
				score = 5 * weight
			}
			if per < 1 {
				score = 2 * weight
			}
			if per >= 1 {
				score = 0 * weight
			}

		}
		scores[k] = score
	}
	return
}

// Node's score depend on node's resources.
func (a *algorithm) _scoreResources() (scores map[string]int, err error) {
	weight := 10
	scores = map[string]int{}
	GB := uint64(1 << 20) // kb

	for name, n := range a.potentialNodes {
		scores[name] = 0
		if val, ok := a.capacityMap[n.Name]; ok {
			if val.CpuUsage > 90 || val.MemoryUsage > 80 {
				continue
			}
			switch {
			case val.CpuUsage > 30 && val.CpuUsage <= 90:
				{
					scores[name] += 1
				}
			case val.CpuUsage <= 30:
				{
					scores[name] += 2
				}
			}
			switch {
			case val.MemoryUsage > 30 && val.MemoryUsage <= 80:
				{
					scores[name] += 1
				}
			case val.MemoryUsage <= 30:
				{
					scores[name] += 2
				}
			}

			if val.Allocatable.Memory > 4*GB {
				scores[name] += 2
			}
		}
	}

	for k, _ := range scores {
		scores[k] = scores[k] * weight
	}
	return
}
