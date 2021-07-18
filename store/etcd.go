package store

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
	"time"
)

type ETCD struct {
	Type string
	Client 	*clientv3.Client
}

func (etcd *ETCD) Get(key string) (resp string, err error)  {
	timeout := time.Second * 10
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	res, err := etcd.Client.Get(ctx,key)
	cancel()
	if err != nil {
		return "", err
	}
	if len(res.Kvs) == 0 {
		return "", err
	}
	return string(res.Kvs[0].Value), err
}
func (etcd *ETCD) Put(key string, val string) (err error) {
	logrus.Debug("etcd start")
	timeout := time.Second * 10
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	logrus.Debug("etcd put")
	_, err = etcd.Client.Put(ctx, key, val)
	if err != nil {
		cancel()
	}
	return
}

func NewETCD(endpoints []string) (etcd *ETCD,err error){
	etcd = &ETCD{}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return
	}
	etcd.Client = cli
	return
}
