package console

import (
	"context"
	"github.com/jin06/binlogo/configs"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	configs.DefaultEnv()
	configs.InitConfigs()
	go func() {
		err := Run(context.Background())
		if err != nil {
			t.Fail()
		}
	}()
	time.Sleep(time.Millisecond * 100)
}
