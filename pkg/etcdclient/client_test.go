package etcdclient

import (
	"github.com/jin06/binlogo/configs"
	"testing"
)

func TestClient(t *testing.T) {
	configs.InitGoTest()
	cli := Default()
	if cli == nil {
		t.Fail()
	}
}
