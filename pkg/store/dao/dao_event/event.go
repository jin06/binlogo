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

func List(res_type string, res_name string, opts ...clientv3.OpOption) (list []*event.Event, err error) {
	key := EventPrefix() + "/" + res_type + "/" + res_name
	resp, err := etcd_client.Default().Get(context.Background(), key, opts...)
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
