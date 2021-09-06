package etcd

import (
	"context"
	"github.com/jin06/binlogo/config"
	"github.com/jin06/binlogo/store/etcd/options"
	"github.com/jin06/binlogo/store/model"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
	"time"
)

var E *ETCD

func DefaultETCD() {
	etcd, err := NewETCD(
		options.Endpoints(config.Cfg.Store.Etcd.Endpoints),
		options.Prefix("binlogo/" + config.Cfg.Cluster.Name),
		options.Timeout(5*time.Second),
	)
	if err != nil {
		panic(err)
	}
	E = etcd
	logrus.Info("default etcd init...")
}

type ETCD struct {
	Client  *clientv3.Client
	options.Options
}

func NewETCD(opt ...options.Option) (etcd *ETCD, err error) {
	ops := options.Options{}
	for _, o := range opt {
		o(&ops)
	}
	etcd = &ETCD{
		Options: ops,
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcd.Options.Endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return
	}
	etcd.Client = cli
	return
}

func (e *ETCD) Read(key string) (resp string, err error) {
	timeout := time.Second * 10
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	res, err := e.Client.Get(ctx, key)
	defer cancel()
	if err != nil {
		return "", err
	}
	if len(res.Kvs) == 0 {
		return "", err
	}
	return string(res.Kvs[0].Value), err
}

func (e *ETCD) Write(key string, val string) (err error) {
	logrus.Debug("etcd start")
	ctx, cancel := context.WithTimeout(context.Background(), e.Timeout)
	logrus.Debug("etcd put")
	_, err = e.Client.Put(ctx, key, val)
	defer cancel()
	return
}

func (e *ETCD) Create(m model.Model) (ok bool, err error) {
	ctx, cancel := context.WithTimeout(context.TODO(), e.Timeout)
	defer cancel()
	key := "/"  + e.Prefix + "/" + m.Key()
	val := m.Val()
	_, err = e.Client.Put(ctx, key, val)
	if err != nil {
		return
	}
	ok = true
	return
}

func (e *ETCD) Update(m model.Model) (ok bool, err error) {
	ctx, cancel := context.WithTimeout(context.TODO(), e.Timeout)
	defer cancel()
	_, err = e.Client.Put(ctx, m.Key(), m.Val())
	if err != nil {
		return
	}
	ok = true
	return
}

func (e *ETCD) Delete(m model.Model) (ok bool, err error) {
	ctx, cancel := context.WithTimeout(context.TODO(), e.Timeout)
	defer cancel()
	_, err = e.Client.Delete(ctx, m.Key())
	if err != nil {
		return
	}
	ok = true
	return
}

func (e *ETCD) Get(m model.Model) (ok bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.Timeout)
	defer cancel()
	res, err := e.Client.Get(ctx, m.Key())
	if err != nil {
		return
	}
	if len(res.Kvs) == 0 {
		return
	}
	if err = m.Unmarshal(res.Kvs[0].Value); err != nil {
		return
	}
	ok = true
	return
}
