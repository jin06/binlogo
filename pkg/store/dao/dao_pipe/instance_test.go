package dao_pipe

import (
	"testing"

	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/pkg/util/random"
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
