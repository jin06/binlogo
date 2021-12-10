package etcd

import (
	"fmt"

	"github.com/jin06/binlogo/configs"
	"github.com/spf13/viper"
)

// Prefix will deprecated
func Prefix() string {
	return fmt.Sprintf("/%s/%s", configs.APP, viper.GetString("cluster.name"))
}
