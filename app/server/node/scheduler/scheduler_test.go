package scheduler

import (
	"context"
	"github.com/jin06/binlogo/configs"
	"testing"
)

func TestRun(t *testing.T) {
	configs.DefaultEnv()
	configs.InitConfigs()
	s := New()
	err := s.Run(context.Background())
	if err != nil {
		t.Fail()
	}
}
