package pipe_tool

import (
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"testing"
)

func TestIsFilterWithName(t *testing.T) {
	filters := []*pipeline.Filter{
		&pipeline.Filter{pipeline.FILTER_BLACK, "mysql"},
		&pipeline.Filter{pipeline.FILTER_BLACK, "base1"},
		&pipeline.Filter{pipeline.FILTER_WHITE, "base1.tbl1"},
	}
	f := NewFilter(filters)
	ok , err := f.IsFilterWithName("mysql.user")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if !ok {
		t.Error("Verify rule error")
		t.Fail()
	}
	ok, err = f.IsFilterWithName("base1.tbl2")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if !ok {
		t.Fail()
	}
	ok, err = f.IsFilterWithName("base1.tbl1")
	if err != nil {
		t.Fail()
	}
	if ok {
		t.Fail()
	}
}
