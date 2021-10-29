package pipeline

import (
	"github.com/jin06/binlogo/pkg/store/dao"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/store/model/scheduler"
)

func CompleteInfo(i *Item, pb *scheduler.PipelineBind) (err error) {
	var nodeName string
	for k,v := range pb.Bindings {
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
	n, err := dao.GetNode(nodeName)
	if err != nil {
		return
	}
	i.Info.BindNode = n
	return
}

func CompleteInfoList(list []*Item, pb *scheduler.PipelineBind) (err error ) {
	for _, v := range list {
		if err = CompleteInfo(v, pb) ; err != nil {
			return
		}
	}
	return
}
