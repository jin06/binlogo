package dao_pipe

import (
	"testing"

	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/util/random"
)

func TestInstance(t *testing.T) {
	configs.InitGoTest()
	pipeName := "go_test_pipeline" + random.String()
	_, err := GetInstance(pipeName)
	if err != nil {
		t.Error(err)
	}
	_, err = AllInstance()
	if err != nil {
		t.Error(err)
	}
	_, err = AllInstanceMap()
	if err != nil {
		t.Error(err)
	}
}
