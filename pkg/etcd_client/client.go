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
		Password: viper.GetString("etcd.password"),
	}
	cli, err = clientv3.New(cfg)
	return
}

var client *clientv3.Client

func Default() *clientv3.Client {
	if client == nil {
		var err error
		client, err = New()
		if err != nil {
			panic(err)
		}
	}
	return client
}
