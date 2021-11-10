package pipeline

import (
	"github.com/jin06/binlogo/pkg/store/dao/dao_node"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
	"github.com/jin06/binlogo/pkg/store/dao/dao_sche"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/store/model/scheduler"
	"github.com/sirupsen/logrus"
)

func CompleteInfo(i *Item) (err error) {
	err = CompletePipelineBind(i, nil)
	if err != nil {
		return
	}
	err = CompletePipelineRun(i,nil)
	if err != nil {
		return
	}
	return
}

func CompletePipelineRun(i *Item, pMap map[string]string) (err error) {
	if pMap != nil {
		i.Info.RunNodeName = pMap[i.Pipeline.Name]
		return
	}
	nodeName , err := dao_pipe.GetInstance(i.Pipeline.Name)
	if err != nil {
		return
	}
	i.Info.RunNodeName = nodeName
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

func CompleteInfoList(list []*Item) (err error) {
	pb, err := dao_sche.GetPipelineBind()
	if err != nil {
		return
	}
	allInstance, err := dao_pipe.AllInstance()
	logrus.Debug(allInstance)
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
