package pipeline

import (
	"context"

	"github.com/jin06/binlogo/v2/pkg/store/dao"
	"github.com/jin06/binlogo/v2/pkg/store/dao/dao_pipe"
	"github.com/jin06/binlogo/v2/pkg/store/dao/dao_sche"
	"github.com/jin06/binlogo/v2/pkg/store/model"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
)

// Item for front display
type Item struct {
	Pipeline *pipeline.Pipeline `json:"pipeline"`
	Info     *Info              `json:"info"`
}

// Info for front display
type Info struct {
	BindNode *node.Node         `json:"bind_node"`
	Instance *pipeline.Instance `json:"instance"`
}

// CompleteInfo complete info for item
func CompleteInfo(i *Item) (err error) {
	err = CompletePipelineBind(i, nil)
	if err != nil {
		return
	}
	err = completePipelineRun(i, nil)
	if err != nil {
		return
	}
	return
}

func completePipelineRun(i *Item, pMap map[string]*pipeline.Instance) (err error) {
	if i == nil {
		return
	}
	if i.Pipeline == nil {
		return
	}
	i.Info.Instance = &pipeline.Instance{}
	if pMap != nil {
		if _, ok := pMap[i.Pipeline.Name]; ok {
			i.Info.Instance = pMap[i.Pipeline.Name]
		}
		return
	}

	ins, err := dao_pipe.GetInstance(i.Pipeline.Name)
	if err != nil {
		return
	}
	i.Info.Instance = ins
	return
}

// CompletePipelineBind complete pipeline bind info
func CompletePipelineBind(i *Item, pb *model.PipelineBind) (err error) {
	if pb == nil {
		pb, err = dao_sche.GetPipelineBind(context.Background())
		if err != nil {
			return
		}
	}
	var nodeName string
	for k, v := range pb.Bindings {
		if k == i.Pipeline.Name {
			nodeName = v
			break
		}
	}
	if i.Info == nil {
		i.Info = &Info{}
	}
	if nodeName == "" {
		i.Info.BindNode = &node.Node{}
		return
	}
	n, err := dao.GetNode(context.Background(), nodeName)
	if err != nil {
		return
	}
	i.Info.BindNode = n
	return
}

func completeEmpty(i *Item) {
	sender := i.Pipeline.Output.Sender
	if sender.RocketMQ == nil {
		sender.RocketMQ = &pipeline.RocketMQ{}
	}
	if sender.Redis == nil {
		sender.Redis = &pipeline.Redis{}
	}
	if sender.Kafka == nil {
		sender.Kafka = &pipeline.Kafka{}
	}
	if sender.Http == nil {
		sender.Http = &pipeline.Http{}
	}
}

// CompleteInfoList complete pipeline list info
func CompleteInfoList(list []*Item) (err error) {
	if len(list) == 0 {
		return nil
	}
	pb, err := dao_sche.GetPipelineBind(context.Background())
	if err != nil {
		return
	}
	allInstance, err := dao_pipe.AllInstanceMap()
	if err != nil {
		return
	}
	for _, v := range list {
		if err1 := CompletePipelineBind(v, pb); err1 != nil {
			continue
		}
		if err2 := completePipelineRun(v, allInstance); err2 != nil {
			return
		}
		completeEmpty(v)
	}
	return
}
