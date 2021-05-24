package store

import (
	"github.com/jin06/binlogo/config"
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
			etcd := new(ETCD)
			cli, err := clientv3.New(clientv3.Config{
				Endpoints: config.Cfg.Store.Etcd.Endpoints,
				DialTimeout: 5 * time.Second,
			})
			if err != nil {

			}
			defer cli.Close()
			DefaultStore = etcd
		}
	default:

	}
}

func Get(key string) (resp string, err error) {
	return DefaultStore.Get(key)
}

func Put(key string, val string) (err error) {
	return DefaultStore.Put(key, val)
}
