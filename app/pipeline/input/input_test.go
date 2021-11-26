package input

import (
	"context"
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/jin06/binlogo/pkg/util/random"
	"strconv"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	configs.InitGoTest()
	serverId := random.Uint32()
	pipeName := "go_test_pipeline" + strconv.Itoa(int(serverId))
	pipeModel := &pipeline.Pipeline{
		Name:      pipeName,
		Status:    pipeline.STATUS_RUN,
		AliasName: pipeName,
		Mysql: &pipeline.Mysql{
			Address:  "127.0.0.1",
			Port:     13308,
			User:     "root",
			Password: "123456",
			ServerId: serverId,
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
	dao_pipe.DeletePipeline(pipeName)
	_, err := dao_pipe.CreatePipeline(pipeModel)
	if err != nil {
		t.Error(err)
	}

	i, err := New(OptionsPipeName(pipeName), OptionsNodeName(configs.NodeName))
	if err != nil {
		t.Error(err)
	}
	ctx, _ := context.WithCancel(context.Background())
	err = i.Run(ctx)
	if err != nil {
		t.Fail()
	}
	dao_pipe.DeletePipeline(pipeName)
}

func TestRunCommon(t *testing.T) {
	configs.InitGoTest()
	serverId := random.Uint32()
	pipeName := "go_test_pipeline" + strconv.Itoa(int(serverId))
	pipeModel := &pipeline.Pipeline{
		Name:      pipeName,
		Status:    pipeline.STATUS_RUN,
		AliasName: pipeName,
		Mysql: &pipeline.Mysql{
			Address:  "127.0.0.1",
			Port:     13308,
			User:     "root",
			Password: "123456",
			ServerId: serverId,
			Flavor:   pipeline.FLAVOR_MYSQL,
			Mode:     pipeline.MODE_POSITION,
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
	dao_pipe.DeletePipeline(pipeName)
	_, err := dao_pipe.CreatePipeline(pipeModel)
	if err != nil {
		t.Error(err)
	}

	i, err := New(OptionsPipeName(pipeName), OptionsNodeName(configs.NodeName))
	if err != nil {
		t.Error(err)
	}
	ctx, _ := context.WithCancel(context.Background())
	err = i.Run(ctx)
	if err != nil {
		t.Fail()
	}
	dao_pipe.DeletePipeline(pipeName)
}
