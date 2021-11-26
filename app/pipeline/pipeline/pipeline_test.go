package pipeline

import (
	"context"
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	configs.InitGoTest()
	pipeName := "go_test_pipeline"
	err := dao_pipe.DeletePipeline(pipeName)
	if err != nil {
		t.Error(err)
	}
	pipeModel := &pipeline.Pipeline{
		Name:      pipeName,
		Status:    pipeline.STATUS_RUN,
		AliasName: pipeName,
		Mysql: &pipeline.Mysql{
			Address:  "127.0.0.1",
			Port:     13308,
			User:     "root",
			Password: "123456",
			ServerId: 12999,
			Flavor:   pipeline.FLAVOR_MYSQL,
			Mode:     pipeline.MODE_GTID,
		},
		Filters: []*pipeline.Filter{
			&pipeline.Filter{
				Type: pipeline.FILTER_BLACK,
				Rule: "mysql",
			},
		},
		Output: &pipeline.Output{Sender: &pipeline.Sender{
			Name:   "",
			Type:   pipeline.SNEDER_TYPE_STDOUT,
			Kafka:  nil,
			Stdout: nil,
			Http:   nil,
		}},
		Replicas:   0,
		CreateTime: time.Time{},
		Remark:     "123",
		IsDelete:   false,
	}
	_, err = dao_pipe.CreatePipeline(pipeModel)
	if err != nil {
		t.Error(err)
	}
	pipe, err := New(
		OptionPipeline(pipeModel),
		OptionPosition(&pipeline.Position{}),
	)
	if err != nil {
		t.Error(err)
	}
	pipe.Run(context.Background())
}
