package input

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
	_, err := dao_pipe.CreatePipeline(pipeModel)
	if err != nil {
		t.Error(err)
	}

	i, err := New(OptionsPipeName(pipeName), OptionsNodeName(configs.NodeName))
	if err != nil {
		t.Error(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	err = i.Run(ctx)
	if err != nil {
		t.Fail()
	}
	cancel()
	dao_pipe.DeletePipeline(pipeName)
}
