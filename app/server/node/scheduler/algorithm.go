package scheduler

import (
	"context"
	"errors"

	"github.com/jin06/binlogo/v2/pkg/store/dao"
	"github.com/jin06/binlogo/v2/pkg/store/model"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
)

type algorithm struct {
	pipeline       *pipeline.Pipeline
	allNodes       map[string]*node.Node
	potentialNodes map[string]*node.Node
	nodesScores    map[string]int
	bestNode       *node.Node
	pb             *model.PipelineBind
	capacityMap    map[string]*node.Capacity
}

func newAlgorithm(opts ...optionAlgorithm) *algorithm {
	a := &algorithm{}
	for _, v := range opts {
		v(a)
	}
	//a.allNodes = map[string]*node.Node{}
	a.potentialNodes = map[string]*node.Node{}
	a.nodesScores = map[string]int{}
	a.bestNode = &node.Node{}
	return a
}

func (a *algorithm) cal() (err error) {
	if a.allNodes == nil {
		a.allNodes, err = dao.AllWorkNodesMap()
		if err != nil {
			return
		}
	}
	if a.pb == nil {
		a.pb, err = dao.GetPipelineBind(context.Background())
		if err != nil {
			return
		}
	}
	if a.capacityMap == nil {
		a.capacityMap, err = dao.CapacityMap(context.Background())
		if err != nil {
			return
		}
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

	for k := range a.potentialNodes {
		if a.nodesScores[k] > score {
			a.bestNode = a.potentialNodes[k]
			score = a.nodesScores[k]
		}
	}
	return
}

func (a *algorithm) calScore() (err error) {

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

	for k := range scores {
		scores[k] = scores[k] * weight
	}
	return
}

type optionAlgorithm func(alg *algorithm)

func withAlgPipe(p *pipeline.Pipeline) optionAlgorithm {
	return func(alg *algorithm) {
		alg.pipeline = p
	}
}

func withAlgAllNodes(allNodes map[string]*node.Node) optionAlgorithm {
	return func(alg *algorithm) {
		alg.allNodes = allNodes
	}
}

func withAlgPipeBind(pb *model.PipelineBind) optionAlgorithm {
	return func(alg *algorithm) {
		alg.pb = pb
	}
}

func withAlgCapMap(cm map[string]*node.Capacity) optionAlgorithm {
	return func(alg *algorithm) {
		alg.capacityMap = cm
	}
}
