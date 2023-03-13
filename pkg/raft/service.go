package raft

import (
	"context"
	"github.com/jin06/binlogo/pkg/proto"
	"github.com/jin06/binlogo/pkg/store/cache"
)

type EntryService struct {
	Cache *cache.CacheManager
}

func (es *EntryService) Get(ctx context.Context, req *proto.GetRequest) (resp *proto.GetResponse, err error) {
	mu := req.Menu
	key := req.Key
	entry, err := es.Cache.GetEntry(mu, key)
	resp = &proto.GetResponse{
		Menu: entry.Menu,
		Key:  entry.Key,
		Val:  string(entry.Value),
	}
	return
}
