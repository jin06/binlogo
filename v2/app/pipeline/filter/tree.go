package filter

import (
	"fmt"
	"strings"

	message2 "github.com/jin06/binlogo/v2/app/pipeline/message"
	"github.com/jin06/binlogo/v2/pkg/pipeline/tool"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
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

func newTree(filters []*pipeline.Filter) (res tree) {
	res = tree{
		DBBlack:    map[string]bool{},
		TableBlack: map[string]bool{},
		DBWhite:    map[string]bool{},
		TableWhite: map[string]bool{},
	}
	if filters == nil {
		return
	}
	for _, v := range filters {
		if !tool.FilterVerify(v) {
			continue
		}
		arr := strings.Split(v.Rule, ".")
		if len(arr) == 1 {
			if v.Type == pipeline.FILTER_BLACK {
				res.DBBlack[v.Rule] = true
			} else {
				res.DBWhite[v.Rule] = true
			}
		} else {
			if v.Type == pipeline.FILTER_BLACK {
				res.TableBlack[v.Rule] = true
			} else {
				res.TableWhite[v.Rule] = true
			}
		}

	}
	return
}
