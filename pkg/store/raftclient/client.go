package raftclient

import (
	"context"
	_ "github.com/Jille/grpc-multi-resolver"
	"github.com/jin06/binlogo/pkg/proto"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/health"
)

type RaftClient struct {
	Target string
	Config string
}

func (r RaftClient) Get(ctx context.Context, in *proto.GetRequest, opts ...grpc.CallOption) (*proto.GetResponse, error) {
	panic("implement me")
}

func (r RaftClient) Set(ctx context.Context, in *proto.SetRequest, opts ...grpc.CallOption) (*proto.SetResponse, error) {
	panic("implement me")
}

func (r RaftClient) Delete(ctx context.Context, in *proto.DelRequest, opts ...grpc.CallOption) (*proto.DelResponse, error) {
	panic("implement me")
}

func NewRaftClient(target string, cfg string) *RaftClient {
	rc := &RaftClient{
		Target: target,
		Config: cfg,
	}
	return rc
}
