package etcd

import (
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"testing"
)

func Init() {
	configs.DefaultEnv()
	configs.InitConfigs()
	DefaultETCD()
}

func TestCreate(t *testing.T) {
	Init()
	testPipe := pipeline.NewPipeline("go_test_pipeline")
	ok, err := Create(testPipe)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Fail()
	}
}

func TestUpdate(t *testing.T) {
	Init()
	testPipe := pipeline.NewPipeline("go_test_pipeline")
	testPipe.Remark = "update"
	ok, err := Update(testPipe)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Fail()
	}
}

func TestDelete(t *testing.T) {
	Init()
	testPipe := pipeline.NewPipeline("go_test_pipeline")
	testPipe.Remark = "update"
	ok , err := Delete(testPipe)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Fail()
	}
}
