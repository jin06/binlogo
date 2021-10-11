package etcd

import (
	"context"
	options2 "github.com/jin06/binlogo/pkg/store/etcd/options"
	model2 "github.com/jin06/binlogo/pkg/store/model"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.etcd.io/etcd/clientv3"
	"os"
	"time"
)

var E *ETCD

func DefaultETCD() {
	//prefix := "binlogo/" + viper.GetString("cluster.name")
	prefix := Prefix()
	etcd, err := NewETCD(
		//options.Endpoints(config.Cfg.Store.Etcd.Endpoints),
		options2.Endpoints(viper.GetStringSlice("store.etcd.endpoints")),
		//options.Prefix("binlogo/"+config.Cfg.Cluster.Name),
		options2.Prefix(prefix),
		options2.Timeout(5*time.Second),
	)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	E = etcd
	logrus.Info("default etcd init...")
}

type ETCD struct {
	Client *clientv3.Client
	options2.Options
}

func NewETCD(opt ...options2.Option) (etcd *ETCD, err error) {
	ops := options2.Options{}
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
func Read(key string) (resp string, err error) {
	return E.Read(key)
}
func (e *ETCD) Write(key string, val string) (err error) {
	logrus.Debug("etcd start")
	ctx, cancel := context.WithTimeout(context.Background(), e.Timeout)
	logrus.Debug("etcd put")
	_, err = e.Client.Put(ctx, key, val)
	defer cancel()
	return
}

func Write(key string, val string) (err error) {
	return E.Write(key, val)
}

func (e *ETCD) Create(m model2.Model, opts ...clientv3.OpOption) (ok bool, err error) {
	ctx, cancel := context.WithTimeout(context.TODO(), e.Timeout)
	defer cancel()
	key := "/" + e.Prefix + "/" + m.Key()
	val := m.Val()
	_, err = e.Client.Put(ctx, key, val, opts...)
	if err != nil {
		return
	}
	ok = true
	return
}

func Create(m model2.Model, opts ...clientv3.OpOption) (bool, error) {
	return E.Create(m, opts...)
}
func (e *ETCD) Update(m model2.Model) (ok bool, err error) {
	ctx, cancel := context.WithTimeout(context.TODO(), e.Timeout)
	defer cancel()
	key := "/" + e.Prefix + "/" + m.Key()
	val := m.Val()
	_, err = e.Client.Put(ctx, key, val)
	if err != nil {
		return
	}
	ok = true
	return
}

func Update(m model2.Model) (bool, error) {
	return E.Update(m)
}

func (e *ETCD) Delete(m model2.Model) (ok bool, err error) {
	ctx, cancel := context.WithTimeout(context.TODO(), e.Timeout)
	defer cancel()
	key := "/" + e.Prefix + "/" + m.Key()
	_, err = e.Client.Delete(ctx, key)
	if err != nil {
		return
	}
	ok = true
	return
}

func Delete(m model2.Model) (bool, error) {
	return E.Delete(m)
}
func (e *ETCD) Get(m model2.Model) (ok bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.Timeout)
	defer cancel()
	key := "/" + e.Prefix + "/" + m.Key()
	res, err := e.Client.Get(ctx, key)
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

func Get(m model2.Model) (bool, error) {
	return E.Get(m)
}