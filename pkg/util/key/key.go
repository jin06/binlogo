package key

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

func GetKey() string {
	return fmt.Sprintf("%s.%d", viper.GetString("node.name"), time.Now().UnixNano())
}
