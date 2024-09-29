package etcd

import (
	"context"
	"fmt"
	"os"
	"time"

	options2 "github.com/jin06/binlogo/v2/pkg/store/etcd/options"
	model2 "github.com/jin06/binlogo/v2/pkg/store/model"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var E *ETCD

// DefaultETCD @Deprecated
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

// NewETCD @Deprecated
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
		Password:    viper.GetString("etcd.password"),
	})
	if err != nil {
		return
	}
	etcd.Client = cli
	return
}

// Read Deprecated
func (e *ETCD) Read(key string) (resp string, err error) {
	timeout := time.Second * 10
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
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

// Read deprecated
func Read(key string) (resp string, err error) {
	return E.Read(key)
}

// Write deprecated
func (e *ETCD) Write(key string, val string) (err error) {
	logrus.Debug("etcd start")
	ctx, cancel := context.WithTimeout(context.TODO(), e.Timeout)
	logrus.Debug("etcd put")
	_, err = e.Client.Put(ctx, key, val)
	defer cancel()
	return
}

// Write deprecated
func Write(key string, val string) (err error) {
	return E.Write(key, val)
}

// Create deprecated
func (e *ETCD) Create(m model2.Model, opts ...clientv3.OpOption) (ok bool, err error) {
	ctx, cancel := context.WithTimeout(context.TODO(), e.Timeout)
	defer cancel()
	key := e.Prefix + "/" + m.Key()
	val := m.Val()
	_, err = e.Client.Put(ctx, key, val, opts...)
	if err != nil {
		return
	}
	ok = true
	return
}

// Create will deprecated
func Create(m model2.Model, opts ...clientv3.OpOption) (bool, error) {
	return E.Create(m, opts...)
}

// Update deprecated
func (e *ETCD) Update(m model2.Model) (ok bool, err error) {
	ctx, cancel := context.WithTimeout(context.TODO(), e.Timeout)
	defer cancel()
	key := e.Prefix + "/" + m.Key()
	val := m.Val()
	_, err = e.Client.Put(ctx, key, val)
	if err != nil {
		return
	}
	ok = true
	return
}

// Update will deprecated
func Update(m model2.Model) (bool, error) {
	return E.Update(m)
}

// Delete deprecated
func (e *ETCD) Delete(m model2.Model) (ok bool, err error) {
	ctx, cancel := context.WithTimeout(context.TODO(), e.Timeout)
	defer cancel()
	key := e.Prefix + "/" + m.Key()
	_, err = e.Client.Delete(ctx, key)
	if err != nil {
		return
	}
	ok = true
	return
}

// Delete will deprecated
func Delete(m model2.Model) (bool, error) {
	return E.Delete(m)
}

// Get deprecated
func (e *ETCD) Get(m model2.Model) (ok bool, err error) {
	ctx, cancel := context.WithTimeout(context.TODO(), e.Timeout)
	defer cancel()
	key := e.Prefix + "/" + m.Key()
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

// List deprecated
func (e *ETCD) List(key string) (list []model2.Model, err error) {
	ctx, cancel := context.WithTimeout(context.TODO(), e.Timeout)
	defer cancel()
	key = e.Prefix + "/" + key
	logrus.Debug("list key: ", key)
	res, err := e.Client.Get(ctx, key, clientv3.WithPrefix(), clientv3.WithFromKey())
	if err != nil {
		return
	}
	//fmt.Println(res)
	//fmt.Println(res.Kvs)
	//fmt.Println(res.Count)

	if len(res.Kvs) == 0 {
		return
	}
	for _, v := range res.Kvs {
		fmt.Println(string(v.Key))
		fmt.Println(string(v.Value))
	}

	return
}

// Get next version will deprecated
func Get(m model2.Model) (bool, error) {
	return E.Get(m)
}
