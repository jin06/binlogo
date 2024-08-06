package configs

import (
	"strconv"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	verArr := strings.Split(VERSITON, ".")
	if len(verArr) != 3 {
		t.Fail()
	}
	for _, v := range verArr {
		n, er := strconv.Atoi(v)
		if er != nil {
			t.Fail()
		}
		if strconv.Itoa(n) != v {
			t.Fail()
		}
	}
}
