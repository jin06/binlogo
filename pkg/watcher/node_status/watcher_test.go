package node_status

import (
	"context"
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/store/dao/dao_node"
	"github.com/jin06/binlogo/pkg/util/random"
	"testing"
	"time"
)

func TestWatcher(t *testing.T) {
	configs.InitGoTest()
	w, err := New(dao_node.StatusPrefix())
	if err != nil {
		t.Error(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	_,  err = w.WatchEtcdList(ctx)
	if err != nil {
		t.Error(err)
	}
	nName := "gotest" + random.String()
	_, err = dao_node.CreateOrUpdateStatus(nName)
	if err != nil {
		t.Error(err)
	}
	dao_node.DeleteStatus(nName)
	time.Sleep(time.Millisecond*100)
	cancel()
}
