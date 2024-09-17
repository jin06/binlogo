package etcdclient

import (
	"testing"

	"github.com/jin06/binlogo/v2/configs"
)

func TestClient(t *testing.T) {
	configs.InitGoTest()
	cli := Default()
	if cli == nil {
		t.Fail()
	}
}
