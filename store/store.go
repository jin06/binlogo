package store

import (
	"github.com/jin06/binlogo/config"
	"github.com/sirupsen/logrus"
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
			etcd, err := NewETCD(config.Cfg.Store.Etcd.Endpoints)
			if err != nil {
				panic(err)
			}
			DefaultStore = etcd
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
