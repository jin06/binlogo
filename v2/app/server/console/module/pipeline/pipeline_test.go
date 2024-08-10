package pipeline

import (
	"testing"

	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/pkg/store/dao/dao_pipe"
	pipeline2 "github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
	"github.com/jin06/binlogo/v2/pkg/util/random"
)

func TestGetItemByName(t *testing.T) {
	configs.InitGoTest()
	pipeName := "go_test_pipeline" + random.String()
	_, err := GetItemByName(pipeName)
	if err != nil {
		t.Fail()
	}
	_, err = dao_pipe.CreatePipeline(pipeline2.NewPipeline(pipeName))
	if err != nil {
		t.Error(err)
	}
	defer func() {
		dao_pipe.DeletePipeline(pipeName)
	}()
	_, err = GetItemByName(pipeName)
	if err != nil {
		t.Error(err)
	}
}

func TestPipeStatus(t *testing.T) {
	pipeName := "go_test_pipeline" + random.String()
	_, err := PipeStatus(pipeName)
	if err != nil {
		t.Fail()
	}
	_, err = dao_pipe.CreatePipeline(pipeline2.NewPipeline(pipeName))
	if err != nil {
		t.Error(err)
	}
	defer func() {
		dao_pipe.DeletePipeline(pipeName)
	}()
	_, err = PipeStatus(pipeName)
	if err != nil {
		t.Fail()
	}
}
