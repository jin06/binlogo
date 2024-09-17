package console

import (
	"context"
	"testing"
	"time"

	"github.com/jin06/binlogo/v2/configs"
)

func TestRun(t *testing.T) {
	configs.InitGoTest()
	go func() {
		err := Run(context.Background())
		if err != nil {
			t.Fail()
		}
	}()
	time.Sleep(time.Millisecond * 100)
}
