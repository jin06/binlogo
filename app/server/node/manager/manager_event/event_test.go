package manager_event

import (
	"context"
	"github.com/jin06/binlogo/configs"
	"testing"
)

func TestRun(t *testing.T) {
	configs.DefaultEnv()
	configs.InitConfigs()
	m := New()
	err := m.Run(context.Background())
	if err != nil {
		t.Fail()
	}
	m.Stop()
}
