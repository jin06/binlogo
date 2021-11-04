package watcher

import (
	"github.com/coreos/etcd/clientv3"
)

type Event struct {
	Event *clientv3.Event
	Data interface{}
}
