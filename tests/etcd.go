package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)
func main()  {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:12379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
	}
	defer cli.Close()
	timeout := time.Second * 10
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	resp, err := cli.Put(ctx, "binlog", "my-binlog-value1")
	resp = resp
	//fmt.Println(resp)
	cancel()
	if err != nil {
		// handle error!
	}
	fmt.Println(1111)
	get, err := cli.Get(ctx, "binlog")
	fmt.Println(2222)
	get = get
	if err != nil {
	}
}
