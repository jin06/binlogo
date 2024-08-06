package util

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func TestStartEnd(t *testing.T) {
	page1, limit1 := 0, 10
	start1, end1 := _StartEnd(page1, limit1)
	if start1 != 0 {
		t.Fail()
	}
	if start1 >= end1 {
		t.Fail()
	}
	page2, limit2 := 1, -10
	start2, end2 := _StartEnd(page2, limit2)
	if start2 < 0 {
		t.Fail()
	}
	if end2 < 0 {
		t.Fail()
	}
	if start2 >= end2 {
		t.Fail()
	}

	s, e := StartEnd(&gin.Context{})
	if s != 0 || e != 10 {
		t.Fail()
	}
}
