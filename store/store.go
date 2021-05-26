package store

import (
	"github.com/jin06/binlogo/config"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
	"time"
)

type Store interface {
	Get(key string) (resp string, err error)
	Put(key string, val string) (err error)
}

var DefaultStore Store

func InitDefault() {
	switch config.Cfg.Store.Type {
	case "etcd":
		{
			etcd := &ETCD{}
			cli, err := clientv3.New(clientv3.Config{
				Endpoints: config.Cfg.Store.Etcd.Endpoints,
				//Endpoints: []string{"localhost:12379"},
				DialTimeout: 5 * time.Second,
			})
			if err != nil {
				panic(err)
			}
			etcd.Client = cli
			DefaultStore = etcd
			//defer cli.Close()
			logrus.Infoln("store is etcd")
		}
	default:
		logrus.Error("Does support the store")
	}
}

func Get(key string) (resp string, err error) {
	return DefaultStore.Get(key)
}

func Put(key string, val string) (err error) {
	return DefaultStore.Put(key, val)
}
