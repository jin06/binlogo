package scheduler_binding

import (
	"context"
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/store/dao/dao_sche"
	"github.com/jin06/binlogo/pkg/util/random"
	"testing"
	"time"
)

func TestWatcher(t *testing.T) {
	configs.InitGoTest()
	w, err := New()
	_, err = w.WatchEtcd(context.Background())
	if err != nil {
		t.Error(err)
	}
	pName := "gotest" + random.String()
	nName := "gotest" + random.String()
	dao_sche.UpdatePipelineBind(pName, nName)
	dao_sche.DeletePipelineBind(pName)
	time.Sleep(time.Millisecond*100)
}
