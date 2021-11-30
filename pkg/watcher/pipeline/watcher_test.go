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
	w, err := New(dao_pipe.PipelinePrefix())
	if err != nil {
		t.Error(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	_, err = w.WatchEtcdList(ctx)
	if err != nil {
		t.Error(err)
	}
	dao_pipe.CreatePipeline(pipeline.NewPipeline(pName))
	dao_pipe.DeletePipeline(pName)
	time.Sleep(time.Millisecond*100)
	cancel()
	time.Sleep(time.Millisecond*100)
}
