package dao_event

import (
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/store/model/event"
	"strconv"
	"testing"
	"time"
)

func TestUpdate(t *testing.T) {
	configs.InitGoTest()
	e := event.NewInfoPipeline("go_test_pipeline", "for go test")
	err := Update(e)
	if err != nil {
		t.Error(err)
	}
}

func TestList(t *testing.T) {
	configs.InitGoTest()
	_, err := List(string(event.PIPELINE), "go_test_pipeline")
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteRange(t *testing.T) {
	configs.InitConfigs()
	deleteTime := time.Now().UnixNano()
	from := configs.NodeName + ".0"
	end := configs.NodeName + "." + strconv.FormatInt(deleteTime, 10)
	prefix := PipelinePrefix() + "/" + "go_test_pipeline" + "/"
	deleted, err := DeleteRange(prefix+from, prefix+end)
	if err != nil {
		t.Error(err)
	}
	t.Log("deleted: ", deleted)
}
