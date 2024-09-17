package tool

import (
	"testing"

	"github.com/jin06/binlogo/v2/app/pipeline/message"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
)

func TestIsFilterWithName(t *testing.T) {
	filters := []*pipeline.Filter{
		{pipeline.FILTER_BLACK, "mysql"},
		{pipeline.FILTER_BLACK, "base1"},
		{pipeline.FILTER_WHITE, "base2"},
		{pipeline.FILTER_WHITE, "base1.tbl1"},
	}
	f := NewFilter(filters)
	ok, err := f.IsFilterWithName("mysql.user")
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("Verify rule error")
	}
	if ok, err = f.IsFilterWithName("base2"); ok {
		t.Fail()
	}
	if ok, err = f.IsFilterWithName("base2.tbl1"); ok {
		t.Fail()
	}

	ok, err = f.IsFilterWithName("base1.tbl2")
	if err != nil {
		t.Error(err)
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
		{pipeline.FILTER_WHITE, "base2"},
		{pipeline.FILTER_WHITE, "base1.tbl1"},
	}
	f := NewFilter(filters)
	msg := message.New()
	msg.Content.Head.Database = "mysql"
	msg.Content.Head.Table = "user"

	if ok := f.IsFilter(msg); !ok {
		t.Fail()
	}

	msg.Content.Head.Database = "base1"
	msg.Content.Head.Table = "tbl1"
	if ok := f.IsFilter(msg); ok {
		t.Fail()
	}
	msg.Content.Head.Database = "base2"
	msg.Content.Head.Table = "tbl1"
	if ok := f.IsFilter(msg); ok {
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
	}); ok {
		t.Fail()
	}
	if ok := FilterVerify(&pipeline.Filter{
		Type: pipeline.FILTER_WHITE,
		Rule: "",
	}); ok {
		t.Fail()
	}
	if ok := FilterVerify(&pipeline.Filter{
		Type: pipeline.FILTER_BLACK,
		Rule: "mysql.user",
	}); !ok {
		t.Fail()
	}
}
