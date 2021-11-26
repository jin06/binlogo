package instance

import (
	"github.com/jin06/binlogo/configs"
	"testing"
)

func TestGetInstanceByName(t *testing.T)  {
	configs.DefaultEnv()
	configs.InitConfigs()
	i, err := GetInstanceByName("go_test_pipeline")
	if err != nil {
		t.Fail()
	}
	t.Log(i)
}
