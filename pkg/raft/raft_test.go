package raft

import (
	"context"
	"testing"
)

func TestNewRaftNode(t *testing.T) {
	ctx := context.Background()
	//New(ctx,"nodeA",39001,"./testdata")
	rn , err := NewRaftNode(ctx, "nodeA", "0.0.0.0", 39001,"./testdata")
	if err != nil {
		t.Error(err)
	}
	_ = rn
}
