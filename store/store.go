package store

import (
	"context"
	"github.com/jin06/binlogo/config"
	log "github.com/sirupsen/logrus"
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
			etcd := ETCD{}
			cli, err := clientv3.New(clientv3.Config{
				//Endpoints: config.Cfg.Store.Etcd.Endpoints,
				Endpoints: []string{"localhost:12379"},
				DialTimeout: 5 * time.Second,
			})
			ctx , _ := context.WithTimeout(context.Background(), 10 * time.Second)
			cli.Put(ctx, "kk", "vv")
			etcd.Client = cli
			if err != nil {
				panic(err)
			}
			defer cli.Close()
			DefaultStore = &etcd
			log.Infoln("store is etcd")
		}
	default:
		log.Error("Does support the store")
	}
}

func Get(key string) (resp string, err error) {
	return DefaultStore.Get(key)
}

func Put(key string, val string) (err error) {
	return DefaultStore.Put(key, val)
}
