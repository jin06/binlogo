package pipe_tool

import (
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"strings"
)

func FilterVerifyStr(s string) bool {
	res := strings.Split(s, ".")
	if len(res) > 2 {
		return false
	}
	if len(res) == 0 {
		return false
	}
	return true
}

func FilterVerify(f *pipeline.Filter) bool {

	if f.Type != pipeline.FILTER_BLACK && f.Type != pipeline.FILTER_WHITE{
		return false
	}
	res := strings.Split(f.Rule, ".")
	if len(res) != 2 && len(res) != 1 {
		return false
	}
	return true
}
