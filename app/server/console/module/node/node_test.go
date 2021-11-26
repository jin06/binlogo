package node

import (
	"github.com/jin06/binlogo/configs"
	"testing"
)

func TestGetAllCount(t *testing.T) {
	configs.DefaultEnv()
	configs.InitConfigs()
	_, err := GetAllCount()
	if err != nil {
		t.Error(err)
	}
}

func TestGetAllRegCount(t *testing.T) {
	configs.DefaultEnv()
	configs.InitConfigs()
	_, err := GetAllRegCount()
	if err != nil {
		t.Fail()
	}
}
