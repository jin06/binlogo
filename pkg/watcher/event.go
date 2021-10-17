package watcher

import (
	model2 "github.com/jin06/binlogo/pkg/store/model"
	"go.etcd.io/etcd/clientv3"
)

type Event struct {
	Event *clientv3.Event
	Model model2.Model
}
