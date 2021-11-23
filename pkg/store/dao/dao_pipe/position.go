package dao_pipe

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/coreos/etcd/clientv3"
	"github.com/jin06/binlogo/pkg/etcd_client"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
)

// PositionPrefix returns etcd prefix of pipeline position
func PositionPrefix() string {
	return etcd_client.Prefix() + "/pipeline/position"
}

// UpdatePosition update pipeline position in etcd
func UpdatePosition(p *pipeline.Position) (err error) {
	key := PositionPrefix() + "/" + p.PipelineName
	b, err := json.Marshal(p)
	if err != nil {
		return
	}
	_, err = etcd_client.Default().Put(context.TODO(), key, string(b))
	return
}

// UpdatePositionSafe  update pipeline position in etcd in safe mode.
// there will be version judgment when updating
func UpdatePositionSafe(pipeName string, opts ...pipeline.OptionPosition) (ok bool, err error) {
	if pipeName == "" {
		err = errors.New("empty pipeline name")
		return
	}
	key := PositionPrefix() + "/" + pipeName
	res, err := etcd_client.Default().Get(context.Background(), key)
	if err != nil {
		return
	}
	revision := int64(0)
	pos := &pipeline.Position{}
	if len(res.Kvs) != 0 {
		revision = res.Kvs[0].CreateRevision
		err = json.Unmarshal(res.Kvs[0].Value, pos)
		if err != nil {
			return
		}
	}
	pos.PipelineName = pipeName
	for _, v := range opts {
		v(pos)
	}
	txn := etcd_client.Default().Txn(context.Background()).
		If(clientv3.Compare(clientv3.CreateRevision(key), "=", revision)).
		Then(clientv3.OpPut(key, pos.Val()))
	resp, err := txn.Commit()
	if err != nil {
		return
	}
	ok = resp.Succeeded
	return
}

// GetPosition get pipeline position from etcd
func GetPosition(pipeName string) (p *pipeline.Position, err error) {
	key := PositionPrefix() + "/" + pipeName
	res, err := etcd_client.Default().Get(context.TODO(), key)
	if err != nil {
		return
	}
	if len(res.Kvs) == 0 {
		return
	}
	p = &pipeline.Position{}
	if err = p.Unmarshal(res.Kvs[0].Value); err != nil {
		return
	}
	return
}

// DeletePosition delete pipeline positon by pipeline name in etcd
func DeletePosition(name string) (err error) {
	if name == "" {
		return errors.New("empty name")
	}
	key := PositionPrefix() + "/" + name
	_, err = etcd_client.Default().Delete(context.TODO(), key)
	return
}
