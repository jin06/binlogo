package etcd_client

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/jin06/binlogo/configs"
	"github.com/spf13/viper"
	"sync"
	"time"
)

// Prefix return cluster prefix
func Prefix() string {
	return fmt.Sprintf("/%s/%s", configs.APP, viper.GetString("cluster.name"))
}

// New return a etcd Client object
func New() (cli *clientv3.Client, err error) {
	cfg := clientv3.Config{
		Endpoints:   viper.GetStringSlice("etcd.endpoints"),
		DialTimeout: 5 * time.Second,
		Password:    viper.GetString("etcd.password"),
	}
	cli, err = clientv3.New(cfg)
	return
}

var client *clientv3.Client

// Default global default clientv3.Client
var mutex sync.Mutex

func Default() *clientv3.Client {
	mutex.Lock()
	defer mutex.Unlock()
	if client == nil {
		var err error
		client, err = New()
		if err != nil {
			panic(err)
		}
	}
	return client
}
