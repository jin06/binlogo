package raftclient

import (
	_ "github.com/Jille/grpc-multi-resolver"
	_ "google.golang.org/grpc/health"
)

type RaftClient struct {
	Target string
	Config string
}

func NewRaftClient(target string, cfg string) *RaftClient {
	rc := &RaftClient{
		Target: target,
		Config: cfg,
	}
	return rc
}
