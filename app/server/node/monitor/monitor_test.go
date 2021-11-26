package monitor

import (
	"context"
	"github.com/jin06/binlogo/configs"
	"testing"
)

func TestRun(t *testing.T) {
	configs.DefaultEnv()
	configs.InitConfigs()
	m, err := NewMonitor()
	if err != nil {
		t.Error(err)
	}
	err = m.Run(context.Background())
	if err != nil {
		t.Error(err)
	}
}
