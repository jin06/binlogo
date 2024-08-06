package pipeline

import "testing"

func TestFilter(t *testing.T) {
	f := BlackFilter("db.tbl")
	if f.Type != FILTER_BLACK {
		t.Fail()
	}
	f = WhiteFilter("db.tbl")
	if f.Type != FILTER_WHITE {
		t.Fail()
	}
}
