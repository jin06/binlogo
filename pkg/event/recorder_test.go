package event

import (
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
