package scheduler

import (
	"errors"
	"github.com/jin06/binlogo/pkg/store/dao"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
)

type algorithm struct {
	pipeline       *pipeline.Pipeline
	allNodes       []*node.Node
	potentialNodes []*node.Node
	nodesScores    map[string]int
	bestNode       *node.Node
}

func newAlgorithm(p *pipeline.Pipeline) *algorithm {
	a := &algorithm{}
	a.pipeline = p
	a.allNodes = []*node.Node{}
	a.potentialNodes = []*node.Node{}
	a.nodesScores = map[string]int{}
	a.bestNode = &node.Node{}
	return a
}

func (a *algorithm) cal() (err error) {
	a.allNodes, err = dao.AllWorkNodes()
	if err != nil {
		return
	}
	err = a.calPotentialNodes()
	if err != nil {
		return err
	}
	err = a.calScores()
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

func (a *algorithm) calScores() (err error) {
	scores := map[string]int{}
	pipeNums := map[string]int{}
	total := 0
	totalPipe := 0
	for _, v := range a.potentialNodes {
		scores[v.Name] = 0
		pipeNums[v.Name] = 0
		total++
	}
	pb, err := dao.GetPipelineBind()
	if err != nil {
		return
	}
	for _, v := range pb.Bindings {
		totalPipe++
		if _, ok := pipeNums[v]; ok {
			pipeNums[v]++
		}
	}
	for k, v := range pipeNums {
		score := 0
		var per float64
		if totalPipe > 0 {
			per = float64(v) / float64(totalPipe) / float64(total)
			if per > 1 {
				score = 0
			}
			if per < 1 {
				score = 2
			}
			if per == 0 {
				score = 5
			}
		}
		scores[k] = score
	}

	a.nodesScores = scores

	return
}

func (a *algorithm) calBestNode() (err error) {
	if len(a.potentialNodes) == 0 {
		err = errors.New("no potential node")
		return
	}
	a.bestNode = a.potentialNodes[0]
	score := a.nodesScores[a.bestNode.Name]

	for k, v := range a.potentialNodes {
		if a.nodesScores[v.Name] > score {
			a.bestNode = a.potentialNodes[k]
			score = a.nodesScores[v.Name]
		}
	}
	return
}
