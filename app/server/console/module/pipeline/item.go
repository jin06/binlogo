package pipeline

import (
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
)

type Item struct {
	Pipeline *pipeline.Pipeline `json:"pipeline"`
	Info *Info
}

type Info struct {
	BindNode *node.Node `json:"bind_node"`
}

