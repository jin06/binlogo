package node

import (
	"context"
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/store/dao/dao_node"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/util/random"
	"testing"
	"time"
)

func TestWatcher(t *testing.T) {
	configs.InitGoTest()
	nName := "gotest" + random.String()
	w, err := New(dao_node.NodePrefix())
	_, err = w.WatchEtcdList(context.Background())
	if err != nil {
		t.Error(err)
	}
	dao_node.CreateNode(node.NewNode(nName))
	dao_node.DeleteNode(nName)

	time.Sleep(time.Millisecond*100)
}
