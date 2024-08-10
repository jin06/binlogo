package dao_pipe

import (
	"testing"
	"time"

	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
	"github.com/jin06/binlogo/v2/pkg/util/random"
)

func TestPipeline(t *testing.T) {
	configs.InitGoTest()
	pipeName := "go_test_pipeline" + random.String()
	pipe := pipeline.Pipeline{
		Name:       pipeName,
		Status:     "",
		AliasName:  "",
		Mysql:      nil,
		Filters:    nil,
		Output:     nil,
		Replicas:   0,
		CreateTime: time.Time{},
		Remark:     "",
		IsDelete:   false,
	}

	ok, err := CreatePipeline(&pipe)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("not ok")
	}
	pipe2, err := GetPipeline(pipeName)
	if err != nil {
		t.Error(err)
	}
	if pipe2.Name != pipeName {
		t.Error("pipe name different")
	}
	ok, err = UpdatePipeline(pipeName,
		pipeline.WithPipeStatus(pipeline.STATUS_RUN),
		pipeline.WithPipeSafe(&pipeline.Pipeline{
			Mysql:     nil,
			AliasName: pipeName,
			Filters:   nil,
			Output:    nil,
			Remark:    pipeName,
		}),
		pipeline.WithPipeDelete(false),
		pipeline.WithPipeMode(pipeline.MODE_POSITION),
		pipeline.WithAddFilter(&pipeline.Filter{}),
		pipeline.WithUpdateFilter(0, &pipeline.Filter{}),
	)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("not update ")
	}
	_, err = AllPipelines()
	if err != nil {
		t.Error(err)
	}
	_, err = AllPipelinesMap()
	if err != nil {
		t.Error(err)
	}
	_, err = DeleteCompletePipeline(pipeName)
	if err != nil {
		t.Error(err)
	}

}
