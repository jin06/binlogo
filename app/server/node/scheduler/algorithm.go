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

func newAlgorithm(p *pipeline.Pipeline) *algorithm{
	a := &algorithm{}
	a.pipeline = p
	a.allNodes = []*node.Node{}
	a.potentialNodes = []*node.Node{}
	a.nodesScores = map[string]int{}
	a.bestNode = &node.Node{}
	return a
}

func (a *algorithm) cal() (err error) {
	a.allNodes, err = dao.AllNodes()
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
	return
}

func (a *algorithm) calBestNode() (err error) {
	if len(a.potentialNodes) > 0 {
		a.bestNode = a.potentialNodes[0]
	} else {
		err = errors.New("no potential node")
	}
	return
}
