package filter

import (
	"testing"

	"github.com/jin06/binlogo/v2/app/pipeline/message"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
)

func TestIsFilter(t *testing.T) {
	tr := &tree{
		DBBlack:    map[string]bool{"mysql": true},
		TableBlack: map[string]bool{"mall.order": true},
		DBWhite:    map[string]bool{"pass": true},
		TableWhite: map[string]bool{"mysql.pass": true},
	}
	testMsg := &message.Message{
		Content: message.Content{
			Head: message.Head{
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
