package watcher

import (
	"github.com/jin06/binlogo/v2/pkg/store/model"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// Event watcher event
type Event struct {
	Event *clientv3.Event
	Data  interface{}
}

type ChangeEvent struct {
	Old any
	New any
}

type EventType byte

func (t EventType) IsDelete() bool {
	return t == EventTypeDelete
}

func (t EventType) IsUpdate() bool {
	return t == EventTypeUpdate
}

const EventTypeUpdate EventType = 0
const EventTypeDelete EventType = 1

type WatcherEvent struct {
	EventType EventType
	Data      model.Model
}
