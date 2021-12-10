package etcdclient

import (
	"fmt"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/jin06/binlogo/configs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
		Username:    viper.GetString("etcd.username"),
	}
	cli, err = clientv3.New(cfg)
	return
}

var client *clientv3.Client

var mutex sync.Mutex

// Default global default clientv3.Client
func Default() *clientv3.Client {
	mutex.Lock()
	defer mutex.Unlock()
	if client == nil {
		var err error
		client, err = New()
		if err != nil {
			//panic(err)
			logrus.Errorln("New default etcd client error: ", err)
		}
	}
	return client
}
