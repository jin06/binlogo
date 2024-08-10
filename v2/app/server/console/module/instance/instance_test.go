package instance

import (
	"testing"
	"time"

	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/pkg/store/dao/dao_pipe"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
	"github.com/jin06/binlogo/v2/pkg/util/random"
)

func TestGetInstanceByName(t *testing.T) {
	configs.InitGoTest()
	pipeName := "go_test_pipeline" + random.String()
	_, err := GetInstanceByName(pipeName)
	if err != nil {
		t.Fail()
	}
	_, err = dao_pipe.CreatePipeline(&pipeline.Pipeline{
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
	})
	if err != nil {
		t.Error(err)
	}
	_, err = GetInstanceByName(pipeName)
	if err != nil {
		t.Error(err)
	}
	_, err = GetInstanceByName("")
	if err == nil {
		t.Fail()
	}

}
