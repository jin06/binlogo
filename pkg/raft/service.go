package raft

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/raft"
	"github.com/jin06/binlogo/pkg/proto"
	"github.com/jin06/binlogo/pkg/store/cache"
	"time"
)

type EntryService struct {
	cache *cache.CacheManager
	raft  *raft.Raft
}

func (es *EntryService) set(entry *cache.Entry) (ok bool, err error) {
	cmd, err := json.Marshal(&entry)
	if err != nil {
		return
	}
	applyFuture := es.raft.Apply(cmd, time.Second)
	err = applyFuture.Error()
	if err != nil {
		return
	}
	ok = true
	return
}

func (es *EntryService) Set(ctx context.Context, request *proto.SetRequest) (response *proto.SetResponse, err error) {
	entry := &cache.Entry{
		Menu:     request.Menu,
		Key:      request.Key,
		Value:    []byte(request.Val),
		Revision: 0,
	}
	ok, err := es.set(entry)
	if err != nil {
		return
	}
	response = &proto.SetResponse{
		Ok: ok,
	}
	return
}

func (es *EntryService) Delete(ctx context.Context, request *proto.DelRequest) (*proto.DelResponse, error) {
	panic("implement me")
}

func (es *EntryService) Get(ctx context.Context, req *proto.GetRequest) (resp *proto.GetResponse, err error) {
	mu := req.Menu
	key := req.Key
	entry, err := es.cache.GetEntry(mu, key)
	resp = &proto.GetResponse{
		Menu: entry.Menu,
		Key:  entry.Key,
		Val:  string(entry.Value),
	}
	return
}
