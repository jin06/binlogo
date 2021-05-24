package store

import (
	"context"
	"go.etcd.io/etcd/clientv3"
)

type ETCD struct {
	Client 	*clientv3.Client
}

func (etcd *ETCD) Get(key string) (resp string, err error)  {
	ctx := context.TODO()
	res, err := etcd.Client.Get(ctx,key)
	if err != nil {
		return "", err
	}
	if len(res.Kvs) == 0 {
		return "", err
	}
	return string(res.Kvs[0].Value), err
}
func (etcd *ETCD) Put (key string, val string) (err error) {
	ctx := context.TODO()
	_, err = etcd.Client.Put(ctx, key, val)
	return err
}