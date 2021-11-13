package filter

import (
	"fmt"
	message2 "github.com/jin06/binlogo/app/pipeline/message"
)

type tree struct {
	DBBlack    map[string]bool
	TableBlack map[string]bool
	DBWhite    map[string]bool
	TableWhite map[string]bool
}

func (t *tree) isFilter(msg *message2.Message) bool {
	if _, ok := t.DBWhite[msg.Content.Head.Database]; ok {
		return false
	}
	table := fmt.Sprintf("%s.%s", msg.Content.Head.Database, msg.Content.Head.Table)
	if _, ok := t.TableWhite[table]; ok {
		return false
	}
	if _, ok := t.DBBlack[msg.Content.Head.Database]; ok {
		return true
	}
	if _, ok := t.TableBlack[table]; ok {
		return true
	}

	return false
}
