package pipeline

import (
	"github.com/jin06/binlogo/pkg/store/dao/dao_node"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
	"github.com/jin06/binlogo/pkg/store/dao/dao_sche"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/jin06/binlogo/pkg/store/model/scheduler"
)

type Item struct {
	Pipeline *pipeline.Pipeline `json:"pipeline"`
	Info     *Info              `json:"info"`
}

type Info struct {
	BindNode *node.Node         `json:"bind_node"`
	Instance *pipeline.Instance `json:"instance"`
	Status   Status             `json:"status"`
}

type Status string

const (
	STATUS_RUN        = "run"
	STATUS_STOP       = "stop"
	STATUS_STOPPING   = "stopping"
	STATUS_SCHEDULING = "scheduling"
	STATUS_SCHEDULED  = "scheduled"
	STATUS_DELETING   = "deleting"
)

func CompleteInfo(i *Item) (err error) {
	err = CompletePipelineBind(i, nil)
	if err != nil {
		return
	}
	err = CompletePipelineRun(i, nil)
	if err != nil {
		return
	}
	return
}

func CompletePipelineRun(i *Item, pMap map[string]*pipeline.Instance) (err error) {
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

func CompletePipelineBind(i *Item, pb *scheduler.PipelineBind) (err error) {
	if pb == nil {
		pb, err = dao_sche.GetPipelineBind()
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
	n, err := dao_node.GetNode(nodeName)
	if err != nil {
		return
	}
	i.Info.BindNode = n
	return
}

func CompletePipelineStatus(i *Item) (err error) {
	if i.Pipeline.IsDelete == true {
		i.Info.Status = STATUS_DELETING
		return
	}
	if i.Pipeline.Status == pipeline.STATUS_STOP {
		if i.Info.BindNode == nil {
			if i.Info.BindNode.Name != "" {
				i.Info.Status = STATUS_STOPPING
				return
			}
		}
		i.Info.Status = STATUS_STOP
		return
	}
	if i.Pipeline.Status == pipeline.STATUS_RUN {
		if i.Info.Instance == nil {
			if i.Info.BindNode == nil {
				i.Info.Status = STATUS_SCHEDULING
				return
			}
			if i.Info.BindNode.Name == "" {
				i.Info.Status = STATUS_SCHEDULING
				return
			}
			i.Info.Status = STATUS_SCHEDULED
			return
		}
	}

	i.Info.Status = STATUS_RUN
	return
}

func CompleteInfoList(list []*Item) (err error) {
	pb, err := dao_sche.GetPipelineBind()
	if err != nil {
		return
	}
	allInstance, err := dao_pipe.AllInstance()
	if err != nil {
		return
	}
	for _, v := range list {
		if err1 := CompletePipelineBind(v, pb); err1 != nil {
			continue
		}
		if err2 := CompletePipelineRun(v, allInstance); err2 != nil {
			return
		}
	}
	return
}
