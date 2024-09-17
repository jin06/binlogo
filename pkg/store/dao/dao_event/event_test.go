package dao_event

import (
	"strconv"
	"testing"
	"time"

	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/pkg/store/model/event"
	clientv3 "go.etcd.io/etcd/client/v3"
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
	_, err := List(string(event.PIPELINE), "go_test_pipeline", 20, clientv3.SortByModRevision, clientv3.SortDescend)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteRange(t *testing.T) {
	configs.InitGoTest()
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
