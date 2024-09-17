package dao_event

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jin06/binlogo/v2/pkg/etcdclient"
	"github.com/jin06/binlogo/v2/pkg/store/model/event"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// EventPrefix returns event prefix
func EventPrefix() string {
	return etcdclient.Prefix() + "/event"
}

// PipelinePrefix returns pipeline event prefix
func PipelinePrefix() string {
	return EventPrefix() + "/pipeline"
}

// NodePrefix returns node event prefix
func NodePrefix() string {
	return EventPrefix() + "/node"
}

func _list(key string, opts ...clientv3.OpOption) (list map[string]*event.Event, err error) {
	opts = append(opts, clientv3.WithPrefix())
	resp, err := etcdclient.Default().Get(context.Background(), key, opts...)
	if err != nil {
		return
	}
	list = map[string]*event.Event{}
	for _, v := range resp.Kvs {
		e := &event.Event{}
		err1 := json.Unmarshal(v.Value, e)
		if err1 != nil {
			logrus.Errorln(err1)
		} else {
			list[string(v.Key)] = e
		}
	}
	return
}

// List returns event list by resource type and resource name
func List(resType string, resName string, limit int64, sortTarget clientv3.SortTarget, sortOrder clientv3.SortOrder) (list map[string]*event.Event, err error) {
	var key string
	if resType == "" || resName == "" {
		err = errors.New("empty resource type or resource name")
		return
	} else {
		key = EventPrefix() + "/" + resType + "/" + resName
	}
	return _list(
		key,
		clientv3.WithLimit(limit),
		clientv3.WithSort(sortTarget, sortOrder),
	)
}

// ScrollList returns event list from key
func ScrollList(key string, n int64, sortTarget clientv3.SortTarget, sortOrder clientv3.SortOrder) (list map[string]*event.Event, err error) {
	if key == "" {
		key = EventPrefix()
	}
	if n <= 0 {
		n = 20
	}
	return _list(
		key,
		clientv3.WithFromKey(),
		clientv3.WithLimit(n),
		clientv3.WithSort(sortTarget, sortOrder),
	)
}

// Update create or update event
func Update(e *event.Event) (err error) {
	if e == nil {
		return errors.New("event is null")
	}
	key := EventPrefix() + "/" + string(e.ResourceType) + "/" + e.ResourceName + "/" + e.Key
	b, err := json.Marshal(e)
	if err != nil {
		return
	}
	_, err = etcdclient.Default().Put(context.Background(), key, string(b))
	return
}

// DeleteRange deletes events
func DeleteRange(fromKey string, endKey string) (deleted int64, err error) {
	res, err := etcdclient.Default().Delete(context.Background(), fromKey, clientv3.WithRange(endKey))
	if err != nil {
		return
	}
	deleted = res.Deleted
	return
}
