package manager_event

import (
	"context"
	"testing"

	"github.com/jin06/binlogo/v2/configs"
)

func TestRun(t *testing.T) {
	configs.InitGoTest()
	m := New()
	err := m.Run(context.Background())
	if err != nil {
		t.Fail()
	}
	err = m.cleanHistoryEvent()
	if err != nil {
		t.Error(err)
	}
	m.Stop()
}
