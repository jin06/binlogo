package key

import (
	"github.com/spf13/viper"
	"testing"
)

func TestGetKey(t *testing.T) {
	viper.SetDefault("node.name", "testkey")
	keys := map[string]bool{}
	for i := 0; i < 10; i++ {
		key := GetKey()
		keys[key] = true
	}
	if len(keys) != 10 {
		t.Fail()
	}
}
