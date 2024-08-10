package tool

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jin06/binlogo/v2/app/pipeline/message"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
)

// NewFilter return a Filter object
func NewFilter(filters []*pipeline.Filter) (f *Filter) {
	f = &Filter{
		DBBlack:    map[string]bool{},
		TableBlack: map[string]bool{},
		DBWhite:    map[string]bool{},
		TableWhite: map[string]bool{},
	}
	for _, v := range filters {
		if !FilterVerify(v) {
			continue
		}
		arr := strings.Split(v.Rule, ".")
		if len(arr) == 1 {
			if v.Type == pipeline.FILTER_BLACK {
				f.DBBlack[v.Rule] = true
			} else {
				f.DBWhite[v.Rule] = true
			}
		} else {
			if v.Type == pipeline.FILTER_BLACK {
				f.TableBlack[v.Rule] = true
			} else {
				f.TableWhite[v.Rule] = true
			}
		}
	}
	return
}

// Filter verify message, filter or pass message by blacklist and whitelist
type Filter struct {
	DBBlack    map[string]bool
	TableBlack map[string]bool
	DBWhite    map[string]bool
	TableWhite map[string]bool
}

// IsFilterWithName filter message by database name and table name.
// return true if not pass
func (t *Filter) IsFilterWithName(name string) (bool, error) {
	res := strings.Split(name, ".")
	lh := len(res)
	switch lh {
	case 1:
		{
			database := res[0]
			if _, ok := t.DBWhite[database]; ok {
				return false, nil
			}
			if _, ok := t.DBBlack[database]; ok {
				return true, nil
			}
		}
	case 2:
		{
			database := res[0]
			if _, ok := t.DBWhite[database]; ok {
				return false, nil
			}
			if _, ok := t.TableWhite[name]; ok {
				return false, nil
			}
			if _, ok := t.DBBlack[database]; ok {
				return true, nil
			}
			if _, ok := t.TableBlack[name]; ok {
				return true, nil
			}
			return false, nil
		}
	default:
		return false, errors.New("wrong rule")
	}
	return false, nil
}

// IsFilter filter message by message object
// return true if not pass
func (t *Filter) IsFilter(msg *message.Message) bool {
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

// FilterVerifyStr verify sting correct
// return false if illegal
func FilterVerifyStr(s string) bool {
	if s == "" {
		return false
	}
	res := strings.Split(s, ".")
	if len(res) > 2 {
		return false
	}
	if len(res) == 0 {
		return false
	}
	return true
}

// FilterVerify verify filter's type
// return false is illegal
func FilterVerify(f *pipeline.Filter) bool {
	if f.Type != pipeline.FILTER_BLACK && f.Type != pipeline.FILTER_WHITE {
		return false
	}
	if f.Rule == "" {
		return false
	}
	res := strings.Split(f.Rule, ".")
	if len(res) != 2 && len(res) != 1 {
		return false
	}
	return true
}
