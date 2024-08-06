package filter

import (
	"testing"

	message2 "github.com/jin06/binlogo/app/pipeline/message"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
)

func TestIsFilter(t *testing.T) {
	tr := &tree{
		DBBlack:    map[string]bool{"mysql": true},
		TableBlack: map[string]bool{"mall.order": true},
		DBWhite:    map[string]bool{"pass": true},
		TableWhite: map[string]bool{"mysql.pass": true},
	}
	testMsg := &message2.Message{
		Content: &message2.Content{
			Head: &message2.Head{
				Database: "mysql",
				Table:    "user",
			},
		},
	}
	if tr.isFilter(testMsg) == false {
		t.Errorf("%s should be filtered. ", testMsg.Table())
		t.Fail()
	}
	testMsg.Content.Head.Database = "mall"
	if tr.isFilter(testMsg) == true {
		t.Errorf("%s should pass ", testMsg.Table())
		t.Fail()
	}
	testMsg.Content.Head.Database = "pass"
	if tr.isFilter(testMsg) == true {
		t.Errorf("%s should pass. ", testMsg.Table())
		t.Fail()
	}
	testMsg.Content.Head.Database = "mysql"
	testMsg.Content.Head.Table = "pass"
	if tr.isFilter(testMsg) == true {
		t.Errorf("%s should pass. ", testMsg.Table())
		t.Fail()
	}
}

func TestNewTree(t *testing.T) {
	filters := []*pipeline.Filter{
		{
			Type: pipeline.FILTER_BLACK,
			Rule: "mysql.user",
		},
	}
	tr := newTree(filters)
	if len(tr.DBBlack) != 0 {
		t.Fail()
	}
	if len(tr.DBWhite) != 0 {
		t.Fail()
	}
	if len(tr.TableBlack) != 1 {
		t.Fail()
	}
	if len(tr.TableWhite) != 0 {
		t.Fail()
	}
}
