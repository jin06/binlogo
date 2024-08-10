package pipeline

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/pkg/store/dao/dao_pipe"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
	"github.com/jin06/binlogo/v2/pkg/util/random"
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
			{
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
	pipe, err := New(
		OptionPipeline(pipeModel),
		OptionPosition(&pipeline.Position{}),
	)
	if err != nil {
		t.Error(err)
	}
	pipe.Run(context.Background())
	time.Sleep(time.Millisecond * 200)
	_, err = dao_pipe.DeletePipeline(pipeName)
	if err != nil {
		t.Error(err)
	}

}

func TestRunCommon(t *testing.T) {
	configs.InitGoTest()
	serverId := random.Uint32()
	pipeName := "go_test_pipeline2" + strconv.Itoa(int(serverId))

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
			{
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
	pipe, err := New(
		OptionPipeline(pipeModel),
		OptionPosition(&pipeline.Position{}),
	)
	if err != nil {
		t.Error(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	pipe.Run(ctx)
	time.Sleep(time.Millisecond * 200)
	_, err = dao_pipe.DeletePipeline(pipeName)
	if err != nil {
		t.Error(err)
	}
	cancel()
}
