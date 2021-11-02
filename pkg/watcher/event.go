package watcher

import (
	"go.etcd.io/etcd/clientv3"
)

type Event struct {
	Event *clientv3.Event
	Data interface{}
}
