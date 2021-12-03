package dao_event

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/coreos/etcd/clientv3"
	"github.com/jin06/binlogo/pkg/etcdclient"
	"github.com/jin06/binlogo/pkg/store/model/event"
	"github.com/sirupsen/logrus"
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

// List returns event list by resource type and resource name
func List(resType string, resName string, opts ...clientv3.OpOption) (list []*event.Event, err error) {
	key := EventPrefix() + "/" + resType + "/" + resName
	opts = append(opts, clientv3.WithPrefix())
	resp, err := etcdclient.Default().Get(context.Background(), key, opts...)
	if err != nil {
		return
	}
	list = []*event.Event{}
	for _, v := range resp.Kvs {
		e := &event.Event{}
		err1 := json.Unmarshal(v.Value, e)
		if err1 != nil {
			logrus.Errorln(err1)
		} else {
			list = append(list, e)
		}
	}
	return
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
