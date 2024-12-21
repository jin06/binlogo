package dao

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jin06/binlogo/v2/pkg/etcdclient"
	"github.com/jin06/binlogo/v2/pkg/store/model"
	storeredis "github.com/jin06/binlogo/v2/pkg/store/redis"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// EventPrefix returns event prefix
func EventPrefix() string {
	return etcdclient.Prefix() + "/event"
}

// PipelinePrefix returns pipeline event prefix
func EventPipelinePrefix() string {
	return EventPrefix() + "/pipeline"
}

// NodePrefix returns node event prefix
func EventNodePrefix() string {
	return EventPrefix() + "/node"
}

func _listEvent(key string, opts ...clientv3.OpOption) (list map[string]*model.Event, err error) {
	opts = append(opts, clientv3.WithPrefix())
	resp, err := etcdclient.Default().Get(context.Background(), key, opts...)
	if err != nil {
		return
	}
	list = map[string]*model.Event{}
	for _, v := range resp.Kvs {
		e := &model.Event{}
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
func ListEvent(resType string, resName string, limit int64, sortTarget clientv3.SortTarget, sortOrder clientv3.SortOrder) (list map[string]*model.Event, err error) {
	var key string
	if resType == "" || resName == "" {
		err = errors.New("empty resource type or resource name")
		return
	} else {
		key = EventPrefix() + "/" + resType + "/" + resName
	}
	return _listEvent(
		key,
		clientv3.WithLimit(limit),
		clientv3.WithSort(sortTarget, sortOrder),
	)
}

// ScrollList returns event list from key
func ScrollListEvent(key string, n int64, sortTarget clientv3.SortTarget, sortOrder clientv3.SortOrder) (list map[string]*model.Event, err error) {
	if key == "" {
		key = EventPrefix()
	}
	if n <= 0 {
		n = 20
	}
	return _listEvent(
		key,
		clientv3.WithFromKey(),
		clientv3.WithLimit(n),
		clientv3.WithSort(sortTarget, sortOrder),
	)
}

// Update create or update event
func UpdateEvent(ctx context.Context, e *model.Event) (bool, error) {
	return storeredis.Default.Update(ctx, e)
}

// DeleteRange deletes events
func DeleteRangeEvent(ctx context.Context, fromKey string, endKey string) (deleted int64, err error) {
	res, err := etcdclient.Default().Delete(context.Background(), fromKey, clientv3.WithRange(endKey))
	if err != nil {
		return
	}
	deleted = res.Deleted
	return
}
