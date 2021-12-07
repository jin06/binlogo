package pipeline

import (
	"context"
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/jin06/binlogo/pkg/util/random"
	"testing"
	"time"
)

func TestWatcher(t *testing.T) {
	configs.InitGoTest()
	pName := "gotest" + random.String()

	WatchList(context.Background(), dao_pipe.PipelinePrefix())
	dao_pipe.CreatePipeline(pipeline.NewPipeline(pName))
	dao_pipe.DeletePipeline(pName)
	time.Sleep(time.Millisecond * 100)
	time.Sleep(time.Millisecond * 100)
}
