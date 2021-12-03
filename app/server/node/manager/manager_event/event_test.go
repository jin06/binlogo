package manager_event

import (
	"context"
	"github.com/jin06/binlogo/configs"
	"testing"
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
