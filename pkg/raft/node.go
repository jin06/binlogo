package raft

import (
	"context"
	"fmt"
	pb "github.com/Jille/raft-grpc-example/proto"
	"github.com/Jille/raft-grpc-leader-rpc/leaderhealth"
	transport "github.com/Jille/raft-grpc-transport"
	"github.com/Jille/raftadmin"
	"github.com/hashicorp/raft"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

// RaftNode a wrapper for raft and grpc server
type RaftNode struct {
	RaftID string
	R      *raft.Raft
	S      *grpc.Server
}

func NewRaftNode(ctx context.Context, raftId string, domain string, port int, dir string) (*RaftNode, error) {
	var err error
	addr := fmt.Sprintf("%s:%d", domain, port)
	fmt.Println(addr)
	fsm := &wordTracker{}
	//todo security
	tm := transport.New(raft.ServerAddress(addr), []grpc.DialOption{grpc.WithInsecure()})
	sock, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	rn := &RaftNode{
		RaftID: raftId,
	}

	rn.R, err = NewRaft(ctx, raftId, addr, fsm, dir, tm.Transport())
	if err != nil {
		return nil, err
	}
	err = rn.bootstrapCluster(addr)
	if err != nil {
		logrus.Errorln(err)
	}
	rn.S = grpc.NewServer()
	pb.RegisterExampleServer(rn.S, &rpcInterface{
		wordTracker: fsm,
		raft:        rn.R,
	})
	tm.Register(rn.S)
	leaderhealth.Setup(rn.R, rn.S, []string{"Example"})
	raftadmin.Register(rn.S, rn.R)
	reflection.Register(rn.S)

	go func(ctx context.Context) {
		if err := rn.S.Serve(sock); err != nil {
			logrus.Fatalln("raft grpc server error, ", err)
		}
	}(ctx)

	return rn, err
}

func (rn *RaftNode) bootstrapCluster(addr string) error {
	cfg := raft.Configuration{
		Servers: []raft.Server{
			{
				Suffrage: raft.Voter,
				ID:       raft.ServerID(rn.RaftID),
				Address:  raft.ServerAddress(addr),
			},
		},
	}
	f := rn.R.BootstrapCluster(cfg)
	if err := f.Error(); err != nil {
		return fmt.Errorf("raft.Raft.BootstrapCluster: %v", err)
	}
	return nil
}
