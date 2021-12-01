package event

import (
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/store/model/event"
	"testing"
	"time"
)

func TestIsExceedTime(t *testing.T) {
	// 2021-11-29 12:00:00
	oldTime := time.Unix(1638158400, 0)
	// 2021-11-29 12:05:00
	newTime := time.Unix(1638158700, 0)
	if isExceedTime(newTime, oldTime) {
		t.Fail()
	}
	// 2021-11-29 12:06:00
	newTime = time.Unix(1638158760, 0)
	if !isExceedTime(newTime, oldTime) {
		t.Fail()
	}
	// 2021-11-29 12:04:00
	newTime = time.Unix(1638158640, 0)
	if isExceedTime(newTime, oldTime) {
		t.Fail()
	}
}

func TestRecorder(t *testing.T) {
	configs.InitGoTest()
	Init()
	Event(event.NewInfoCluster("test message"))
	EventErrorPipeline("go_test_pipeline", "test message")
	EventInfoPipeline("go_test_pipeline", "test message")
	time.Sleep(time.Millisecond * 100)
}
