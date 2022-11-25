package raft

import (
	"context"
	"flag"
	"fmt"
	"github.com/hashicorp/raft"
	"os"
	"path/filepath"

	boltdb "github.com/hashicorp/raft-boltdb"
)

var (
	raftBootstrap = flag.Bool("raft_bootstrap", false, "Whether to bootstrap the Raft cluster")
)

func NewRaft(ctx context.Context, myID, addr string, fsm raft.FSM, raftDir string, trans raft.Transport) (*raft.Raft,  error) {
	c := raft.DefaultConfig()
	c.LocalID = raft.ServerID(myID)

	baseDir := filepath.Join(raftDir, myID)
	err := os.MkdirAll(baseDir,0731)
	if err != nil {
		return nil, fmt.Errorf("create directory error, %v", err)
	}

	ldb, err := boltdb.NewBoltStore(filepath.Join(baseDir, "logs.dat"))
	if err != nil {
		return nil, fmt.Errorf(`boltdb.NewBoltStore(%q): %v`, filepath.Join(baseDir, "logs.dat"), err)
	}

	sdb, err := boltdb.NewBoltStore(filepath.Join(baseDir, "stable.dat"))
	if err != nil {
		return nil, fmt.Errorf(`boltdb.NewBoltStore(%q): %v`, filepath.Join(baseDir, "stable.dat"), err)
	}

	fss, err := raft.NewFileSnapshotStore(baseDir, 3, os.Stderr)
	if err != nil {
		return nil, fmt.Errorf(`raft.NewFileSnapshotStore(%q, ...): %v`, baseDir, err)
	}

	//tm := transport.New(raft.ServerAddress(myAddress), []grpc.DialOption{grpc.WithInsecure()})

	//r, err := raft.NewRaft(c, fsm, ldb, sdb, fss, tm.Transport())
	r, err := raft.NewRaft(c, fsm, ldb, sdb, fss,trans)
	if err != nil {
		return nil, fmt.Errorf("raft.NewRaft: %v", err)
	}

	if *raftBootstrap {
		cfg := raft.Configuration{
			Servers: []raft.Server{
				{
					Suffrage: raft.Voter,
					ID:       raft.ServerID(myID),
					Address:  raft.ServerAddress(addr),
				},
			},
		}
		f := r.BootstrapCluster(cfg)
		if err := f.Error(); err != nil {
			return nil, fmt.Errorf("raft.Raft.BootstrapCluster: %v", err)
		}
	}

	return r, nil
}
