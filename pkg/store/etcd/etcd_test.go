package etcd

import (
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
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

func TestGet(t *testing.T) {
	file := "../../../configs/binlogo.yaml"
	configs.InitViperFromFile(file)
	DefaultETCD()
	m := pipeline.NewPipelineH()
	m.Name = "test"
	_, err := E.GetH(m)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(m)
	t.Log(m.Header.Revision)
}

func TestGet2(t *testing.T) {

}
