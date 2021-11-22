package dao_event

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/coreos/etcd/clientv3"
	"github.com/jin06/binlogo/pkg/etcd_client"
	"github.com/jin06/binlogo/pkg/store/model/event"
	"github.com/sirupsen/logrus"
)

func EventPrefix() string {
	return etcd_client.Prefix() + "/event"
}

func PipelinePrefix() string {
	return EventPrefix() + "/pipeline"
}

func NodePrefix() string  {
	return EventPrefix() + "/node"
}

func List(resType string, resName string, opts ...clientv3.OpOption) (list []*event.Event, err error) {
	key := EventPrefix() + "/" + resType + "/" + resName
	opts = append(opts, clientv3.WithPrefix())
	resp, err := etcd_client.Default().Get(context.Background(), key, opts...)
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

func Update(e *event.Event) (err error) {
	if e == nil {
		return errors.New("event is null")
	}
	key := EventPrefix() + "/" + string(e.ResourceType) + "/" + e.ResourceName + "/" + e.Key
	b, err := json.Marshal(e)
	if err != nil {
		return
	}
	_, err = etcd_client.Default().Put(context.Background(), key, string(b))
	return
}

func DeleteRange(fromKey string, endKey string) (deleted int64, err error) {
	res, err := etcd_client.Default().Delete(context.Background(), fromKey, clientv3.WithRange(endKey), clientv3.WithFromKey())
	if err != nil {
		return
	}
	deleted = res.Deleted
	return
}
