package key

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

// GetKey returns a etcd key
// the key is node name and unix nano time
func GetKey() string {
	return fmt.Sprintf("%s.%d", viper.GetString("node.name"), time.Now().UnixNano())
}
