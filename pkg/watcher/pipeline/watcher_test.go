package pipeline

import (
	"context"
	"github.com/jin06/binlogo/configs"
	etcd2 "github.com/jin06/binlogo/pkg/store/etcd"
	"testing"
)

func TestWatch(t *testing.T) {
	file := "../../../configs/binlogo.yaml"
	configs.InitViperFromFile(file)
	etcd2.DefaultETCD()
	t.Log("start test watch")
	w, err := New("/binlogo/cluster1/pipeline/test")
	if err != nil {
		t.Fatal(err)
	}
	w.Watch(context.Background())
	for p := range w.Queue {
		t.Log(p)
	}
}
