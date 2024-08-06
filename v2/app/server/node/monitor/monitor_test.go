package monitor

import (
	"context"
	"testing"
	"time"

	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/register"
	"github.com/jin06/binlogo/pkg/store/dao/dao_node"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
	"github.com/jin06/binlogo/pkg/store/dao/dao_sche"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/jin06/binlogo/pkg/util/random"
)

func TestRun(t *testing.T) {
	configs.InitGoTest()
	m, err := NewMonitor()
	if err != nil {
		t.Error(err)
	}
	nodeName := "gotest" + random.String()
	nModel := node.NewNode(nodeName)
	dao_node.CreateNode(nModel)
	defer dao_node.DeleteNode(nodeName)
	pipeName := "gotest" + random.String()
	pModel := pipeline.NewPipeline(pipeName)
	pModel.Status = pipeline.STATUS_RUN
	dao_pipe.CreatePipeline(pModel)
	defer dao_pipe.DeletePipeline(pipeName)

	ctx, cancel := context.WithCancel(context.Background())
	err = m.Run(ctx)
	if err != nil {
		t.Error(err)
	}
	dao_sche.UpdatePipelineBind(pipeName, nodeName)
	defer dao_sche.DeletePipelineBind(pipeName)
	r := register.New(
		register.WithTTL(5),
		register.WithKey(dao_node.NodeRegisterPrefix()+"/"+nodeName),
	)
	r.Run(ctx)
	err = m.checkAllNode()
	if err != nil {
		t.Error(err)
	}
	err = checkAllNodeStatus()
	if err != nil {
		t.Error(err)
	}
	err = m.checkAllNodeBind()
	if err != nil {
		t.Error(err)
	}
	m.checkAllPipelineBind()
	err = m.deletePipeline(pModel)
	if err != nil {
		t.Error(err)
	}

	m.checkAllPipelineDelete()
	time.Sleep(time.Millisecond * 500)
	m.Stop(context.Background())

	cancel()
}
