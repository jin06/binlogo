package node

import (
	"testing"

	"github.com/jin06/binlogo/configs"
)

func TestGetAllCount(t *testing.T) {
	configs.InitGoTest()
	_, err := GetAllCount()
	if err != nil {
		t.Error(err)
	}
}

func TestGetAllRegCount(t *testing.T) {
	configs.InitGoTest()
	_, err := GetAllRegCount()
	if err != nil {
		t.Fail()
	}
}
