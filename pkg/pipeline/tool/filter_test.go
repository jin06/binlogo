package tool

import (
	"github.com/jin06/binlogo/app/pipeline/message"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"testing"
)

func TestIsFilterWithName(t *testing.T) {
	filters := []*pipeline.Filter{
		{pipeline.FILTER_BLACK, "mysql"},
		{pipeline.FILTER_BLACK, "base1"},
		{pipeline.FILTER_WHITE, "base1.tbl1"},
	}
	f := NewFilter(filters)
	ok, err := f.IsFilterWithName("mysql.user")
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

func TestIsFilter(t *testing.T) {
	filters := []*pipeline.Filter{
		{pipeline.FILTER_BLACK, "mysql"},
		{pipeline.FILTER_BLACK, "base1"},
		{pipeline.FILTER_WHITE, "base1.tbl1"},
	}
	f := NewFilter(filters)
	msg1 := message.New()
	msg1.Content.Head.Database = "mysql"
	msg1.Content.Head.Table = "user"

	msg2 := message.New()
	msg2.Content.Head.Database = "base1.tbl1"
	msg2.Content.Head.Table = "base1.tbl1"

	if ok := f.IsFilter(msg1); !ok {
		t.Fail()
	}
	if ok := f.IsFilter(msg2); ok {
		t.Fail()
	}
}

func TestFilterVerifyStr(t *testing.T) {
	test1 := "mysql.user"
	if ok := FilterVerifyStr(test1); !ok {
		t.Fail()
	}
	if ok := FilterVerifyStr("mysql.user.1"); ok {
		t.Fail()
	}
	if ok := FilterVerifyStr("mysql"); !ok {
		t.Fail()
	}
	if ok := FilterVerifyStr(""); ok {
		t.Fail()
	}
}

func TestFilterVerify(t *testing.T) {
	if ok := FilterVerify(&pipeline.Filter{
		Type: pipeline.FILTER_BLACK,
		Rule: "a.b.c",
	});ok {
		t.Fail()
	}
	if ok := FilterVerify(&pipeline.Filter{
		Type: pipeline.FILTER_WHITE,
		Rule: "",
	});ok {
		t.Fail()
	}
}
