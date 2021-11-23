package watcher

import (
	"github.com/coreos/etcd/clientv3"
)

// Event watcher event
type Event struct {
	Event *clientv3.Event
	Data  interface{}
}
