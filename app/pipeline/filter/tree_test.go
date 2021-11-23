package filter

import (
	message2 "github.com/jin06/binlogo/app/pipeline/message"
	"testing"
)

func TestisFilter(t *testing.T) {
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
	if !tr.isFilter(testMsg) {
		t.Fail()
	}
	testMsg.Content.Head.Database = "mall"
	if !tr.isFilter(testMsg) {
		t.Fail()
	}
	testMsg.Content.Head.Database = "pass"
	if tr.isFilter(testMsg) {
		t.Fail()
	}
	testMsg.Content.Head.Database = "mysql"
	testMsg.Content.Head.Table = "pass"
	if !tr.isFilter(testMsg) {
		t.Fail()
	}
}
