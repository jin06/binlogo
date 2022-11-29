package raft

import (
	"context"
	"fmt"
	"net"

	pb "github.com/Jille/raft-grpc-example/proto"
	"github.com/Jille/raft-grpc-leader-rpc/leaderhealth"
	transport "github.com/Jille/raft-grpc-transport"
	"github.com/Jille/raftadmin"
	"github.com/hashicorp/raft"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// RaftNode a wrapper for raft and grpc server
type RaftNode struct {
	RaftID      string
	R           *raft.Raft
	S           *grpc.Server
	RaftServers []raft.Server
}

type NodeConfig struct {
	ID      raft.ServerID
	Address raft.ServerAddress
}

func NewRaftNode(ctx context.Context, myServer raft.Server, dir string, bootstrap bool, raftServers []raft.Server) (*RaftNode, error) {
	var err error
	addr := myServer.Address
	fmt.Println(addr)
	fsm := &wordTracker{}
	//todo security
	tm := transport.New(raft.ServerAddress(addr), []grpc.DialOption{grpc.WithInsecure()})
	sock, err := net.Listen("tcp", string(addr))
	if err != nil {
		return nil, err
	}

	rn := &RaftNode{
		RaftID:      string(myServer.ID),
		RaftServers: raftServers,
	}

	rn.R, err = NewRaft(ctx, string(myServer.ID), string(addr), fsm, dir, tm.Transport())
	if err != nil {
		return nil, err
	}
	if bootstrap {
		err = rn.bootstrapCluster(string(addr))

		if err != nil {
			logrus.Errorln(err)
		}
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
	fmt.Println(rn.R.GetConfiguration().Configuration())
	go rn.doLeader(ctx)

	return rn, err
}

func (rn *RaftNode) doLeader(ctx context.Context) {
	ch := rn.R.LeaderCh()
	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		case <-ch:
			{
				rn.initNodes()
			}
		}
	}
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

func (rn *RaftNode) initNodes() {
	m := rn.existiServers()
	for _, v := range rn.RaftServers {
		if _, ok := m[v.ID]; !ok {
			rn.R.AddVoter(v.ID, v.Address, 0, 0)
		}
	}
}

func (rn *RaftNode) existiServers() map[raft.ServerID]raft.Server {
	m := map[raft.ServerID]raft.Server{}
	servers := rn.R.GetConfiguration().Configuration().Servers
	for _, v := range servers {
		m[v.ID] = v
	}
	return m
}
