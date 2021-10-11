package scheduler

import (
	"github.com/jin06/binlogo/pkg/store/model"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
)

type algorithm struct {
	pipeline       *pipeline.Pipeline
	allNodes       []*model.Node
	potentialNodes []*model.Node
	nodesScores    map[string]int
	bestNode       *model.Node
}

func (a *algorithm) cal() (err error) {
	err = a.calPotentialNodes()
	if err != nil {
		return err
	}
	a.calScores()
	a.calBestNode()
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
	}
	return
}
