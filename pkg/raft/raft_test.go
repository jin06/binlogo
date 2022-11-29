package raft

import (
	"context"
	"testing"

	"github.com/hashicorp/raft"
)

func TestNewRaftNode(t *testing.T) {
	ctx := context.Background()
	//New(ctx,"nodeA",39001,"./testdata")
	nodes := []raft.Server{
		raft.Server{raft.Voter, raft.ServerID("nodeA"), raft.ServerAddress("0.0.0.0:39001")},
		raft.Server{raft.Voter, raft.ServerID("nodeB"), raft.ServerAddress("0.0.0.0:39002")},
		raft.Server{raft.Voter, raft.ServerID("nodeC"), raft.ServerAddress("0.0.0.0:39003")},
	}
	rn, err := NewRaftNode(ctx, "nodeA", "0.0.0.0", 39001, "./testdata", false, nodes)
	if err != nil {
		t.Error(err)
	}
	_ = rn
	select {}
}

func TestNewRaftNode2(t *testing.T) {
	ctx := context.Background()
	//New(ctx,"nodeA",39001,"./testdata2")
	nodes := []raft.Server{
		raft.Server{raft.Voter, raft.ServerID("nodeA"), raft.ServerAddress("0.0.0.0:39001")},
		raft.Server{raft.Voter, raft.ServerID("nodeB"), raft.ServerAddress("0.0.0.0:39002")},
		raft.Server{raft.Voter, raft.ServerID("nodeC"), raft.ServerAddress("0.0.0.0:39003")},
	}
	rn, err := NewRaftNode(ctx, "nodeB", "0.0.0.0", 39002, "./testdata", false, nodes)
	if err != nil {
		t.Error(err)
	}
	_ = rn
	select {}
}

func TestNewRaftNode3(t *testing.T) {
	ctx := context.Background()
	//New(ctx,"nodeA",39001,"./testdata2")
	nodes := []raft.Server{
		raft.Server{raft.Voter, raft.ServerID("nodeA"), raft.ServerAddress("0.0.0.0:39001")},
		raft.Server{raft.Voter, raft.ServerID("nodeB"), raft.ServerAddress("0.0.0.0:39002")},
		raft.Server{raft.Voter, raft.ServerID("nodeC"), raft.ServerAddress("0.0.0.0:39003")},
	}
	rn, err := NewRaftNode(ctx, "nodeC", "0.0.0.0", 39003, "./testdata", false, nodes)
	if err != nil {
		t.Error(err)
	}
	_ = rn
	select {}
}
