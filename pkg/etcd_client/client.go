package etcd_client

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/jin06/binlogo/configs"
	"github.com/spf13/viper"
	"time"
)

func Prefix() string {
	return fmt.Sprintf("/%s/%s", configs.APP, viper.GetString("cluster.name"))
}

func New() (cli *clientv3.Client, err error) {
	cfg := clientv3.Config{
		Endpoints:   viper.GetStringSlice("etcd.endpoints"),
		DialTimeout: 5 * time.Second,
	}
	cli, err = clientv3.New(cfg)
	return
}
