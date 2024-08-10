package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/pkg/store/dao/dao_node"
	"github.com/jin06/binlogo/v2/pkg/store/dao/dao_pipe"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
	"github.com/jin06/binlogo/v2/pkg/util/random"
)

func TestRun(t *testing.T) {
	configs.InitGoTest()
	s := New()
	err := s.Run(context.Background())
	if err != nil {
		t.Fail()
	}
	pName := random.String()
	nName := random.String()
	pModel := pipeline.NewPipeline(pName)
	pModel.Status = pipeline.STATUS_RUN
	dao_pipe.CreatePipeline(pModel)
	dao_node.CreateNode(node.NewNode(nName))
	defer dao_pipe.DeletePipeline(pName)
	defer dao_node.DeleteNode(nName)
	err = scheduleOne(pModel)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Millisecond * 200)
	s.Stop()
}
