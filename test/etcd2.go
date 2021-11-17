package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
)

func main() {
	cli, _ := clientv3.New(clientv3.Config{Endpoints: []string{"localhost:12379"}})
	fmt.Println(cli.Leases(context.Background()))
}
