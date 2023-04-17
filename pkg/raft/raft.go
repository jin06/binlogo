package raft

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/raft"

	boltdb "github.com/hashicorp/raft-boltdb"
)

func NewRaft(c *raft.Config, fsm raft.FSM, baseDir string, trans raft.Transport) (*raft.Raft, error) {
	err := os.MkdirAll(baseDir, 0731)
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

	r, err := raft.NewRaft(c, fsm, ldb, sdb, fss, trans)

	//tch := make(chan raft.Observation, 10000)
	//testOb := raft.NewObserver(tch, false, nil)
	//r.RegisterObserver(testOb)
	//go func() {
	//	for v := range tch {
	//		fmt.Println(v)
	//	}
	//}()

	if err != nil {
		return nil, fmt.Errorf("raft.NewRaft: %v", err)
	}

	return r, nil
}
