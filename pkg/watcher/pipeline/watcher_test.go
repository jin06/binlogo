package pipeline

import (
	"context"
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
	etcd2 "github.com/jin06/binlogo/pkg/store/etcd"
	"testing"
)

func TestWatch(t *testing.T) {
	file := "../../../configs/binlogo.yaml"
	configs.InitViperFromFile(file)
	etcd2.DefaultETCD()
	t.Log("start test watch")
	w, err := New(dao_pipe.PipelinePrefix() + "/test")
	if err != nil {
		t.Fatal(err)
	}
	w.Watch(context.Background())
	for p := range w.Queue {
		t.Log(p)
	}
}
