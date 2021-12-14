package watcher

import (
	clientv3 "go.etcd.io/etcd/client/v3"
)

// Event watcher event
type Event struct {
	Event *clientv3.Event
	Data  interface{}
}
