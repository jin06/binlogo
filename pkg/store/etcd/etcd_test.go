package etcd

import (
	"github.com/jin06/binlogo/configs"
	"testing"
)

func TestList(t *testing.T) {
	file := "../../../configs/binlogo.yaml"
	configs.InitViperFromFile(file)
	DefaultETCD()
	res, err := E.List("pipeline")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}
