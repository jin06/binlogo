package pipeline

import (
	"github.com/jin06/binlogo/configs"
	"testing"
)

func TestGetItemByName(t *testing.T) {
	configs.DefaultEnv()
	configs.InitConfigs()
	_, err := GetItemByName("go_test_pipeline")
	if err != nil {
		t.Fail()
	}
}

func TestPipeStatus(t *testing.T) {
	_, err := PipeStatus("go_test_pipeline")
	if err != nil {
		t.Fail()
	}
}
