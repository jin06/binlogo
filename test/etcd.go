package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	model2 "github.com/jin06/binlogo/pkg/store/model"
	"time"
)

func main() {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:12379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		fmt.Println(err)
	}
	defer cli.Close()
	timeout := time.Second * 10
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	database := new(model2.Database)
	database.Name = "testdata"
	js, err := json.Marshal(database)
	if err != nil {
		fmt.Println(err)
	}
	key := "/binlogo/culster1/database/1"
	resp, err := cli.Put(ctx, key, string(js))
	resp = resp
	//fmt.Println(resp)
	if err != nil {
		// handle error!
	}
	fmt.Println(1111)
	get, err := cli.Get(ctx, key)
	fmt.Println(get.Kvs)
	_ = get
	fmt.Println(2222)
	if err != nil {
	}
	cancel()
}
