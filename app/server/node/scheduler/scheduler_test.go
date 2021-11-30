package scheduler

import (
	"context"
	"github.com/jin06/binlogo/configs"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	configs.DefaultEnv()
	configs.InitConfigs()
	s := New()
	err := s.Run(context.Background())
	if err != nil {
		t.Fail()
	}
	time.Sleep(time.Millisecond*200)
	s.Stop()
}
