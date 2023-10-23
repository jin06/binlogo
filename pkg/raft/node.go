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
	RaftID  string
	R       *raft.Raft
	S       *grpc.Server
	RConfig *raft.Config
	*RServers
}

type RServers struct {
	Arr []*raft.Server
	Map map[raft.ServerID]*raft.Server
}

type NodeConfig struct {
	ID      raft.ServerID
	Address raft.ServerAddress
}

func newRServers(srvs []*raft.Server) *RServers {
	rsServers := &RServers{
		Arr: srvs,
	}
	m := map[raft.ServerID]*raft.Server{}
	for _, v := range srvs {
		m[v.ID] = v
	}
	return rsServers
}

func NewRaftNode(ctx context.Context, myServer raft.Server, dir string, raftServers []*raft.Server) (*RaftNode, error) {
	logrus.Infoln("start new raft node", myServer)
	var err error
	addr := myServer.Address
	fsm := &wordTracker{}
	//todo security
	tm := transport.New(myServer.Address, []grpc.DialOption{grpc.WithInsecure()})
	sock, err := net.Listen("tcp", string(addr))
	if err != nil {
		return nil, err
	}

	c := raft.DefaultConfig()
	c.LocalID = myServer.ID

	rn := &RaftNode{
		RaftID:   string(myServer.ID),
		RServers: newRServers(raftServers),
		RConfig:  c,
	}

	rn.R, err = NewRaft(c, fsm, dir, tm.Transport())
	if err != nil {
		return nil, err
	}
	err = rn.bootstrapCluster(string(addr))

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
				logrus.Infoln("win election")
				rn.initNodes()
			}
		}
	}
}

func (rn *RaftNode) bootstrapCluster(addr string) error {
	if len(rn.R.GetConfiguration().Configuration().Servers) > 0 {
		return nil
	}
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
	m := rn.existServers()
	for _, v := range rn.Arr {
		if _, ok := m[v.ID]; !ok {
			if v.Suffrage == raft.Voter {
				rn.R.AddVoter(v.ID, v.Address, 0, 0)
			} else {
				rn.R.AddNonvoter(v.ID, v.Address, 0, 0)
			}
		}
	}
	for _, v := range m {
		if _, ok := rn.Map[v.ID]; !ok {
			rn.R.RemoveServer(v.ID, 0, 0)
		}
	}

}

func (rn *RaftNode) existServers() map[raft.ServerID]raft.Server {
	m := map[raft.ServerID]raft.Server{}
	servers := rn.R.GetConfiguration().Configuration().Servers
	for _, v := range servers {
		m[v.ID] = v
	}
	return m
}
