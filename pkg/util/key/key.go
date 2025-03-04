package key

import (
	"fmt"
	"time"

	"github.com/jin06/binlogo/v2/configs"
)

// GetKey returns a etcd key
// the key is node name and unix nano time
func GetKey() string {
	return fmt.Sprintf("%s.%d", configs.Default.NodeName, time.Now().UnixNano())
}
