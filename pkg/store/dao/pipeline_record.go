package dao

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jin06/binlogo/v2/pkg/etcdclient"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// RecordPrefix returns etcd record prefix
func RecordPrefix() string {
	return etcdclient.Prefix() + "/pipeline/record"
}

// UpdatePosition update pipeline record position in etcd
func UpdateRecord(p *pipeline.RecordPosition) (err error) {
	key := RecordPrefix() + "/" + p.PipelineName
	b, err := json.Marshal(p)
	if err != nil {
		return
	}
	_, err = etcdclient.Default().Put(context.Background(), key, string(b))
	return
}

// UpdateRecordSafe  update pipeline record position in etcd in safe mode.
// there will be version judgment when updating
func UpdateRecordSafe(pipeName string, opts ...pipeline.OptionRecord) (ok bool, err error) {
	if pipeName == "" {
		err = errors.New("empty pipeline name")
		return
	}
	key := RecordPrefix() + "/" + pipeName
	res, err := etcdclient.Default().Get(context.Background(), key)
	if err != nil {
		return
	}
	revision := int64(0)
	record := &pipeline.RecordPosition{}
	if len(res.Kvs) != 0 {
		revision = res.Kvs[0].CreateRevision
		err = json.Unmarshal(res.Kvs[0].Value, record)
		if err != nil {
			return
		}
	}
	record.PipelineName = pipeName
	for _, v := range opts {
		v(record)
	}
	b, err := json.Marshal(record)
	if err != nil {
		return
	}
	txn := etcdclient.Default().Txn(context.Background()).
		If(clientv3.Compare(clientv3.CreateRevision(key), "=", revision)).
		Then(clientv3.OpPut(key, string(b)))
	resp, err := txn.Commit()
	if err != nil {
		return
	}
	ok = resp.Succeeded
	return
}

// GetRecord get pipeline record position from etcd
func GetRecord(pipeName string) (r *pipeline.RecordPosition, err error) {
	key := RecordPrefix() + "/" + pipeName
	res, err := etcdclient.Default().Get(context.Background(), key)
	if err != nil {
		return
	}
	if len(res.Kvs) == 0 {
		return
	}
	r = &pipeline.RecordPosition{}
	if err = json.Unmarshal(res.Kvs[0].Value, r); err != nil {
		return
	}
	return
}
