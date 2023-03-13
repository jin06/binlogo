package raft

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/raft"
	"github.com/jin06/binlogo/pkg/store/cache"
	"github.com/sirupsen/logrus"
	"io"
)

type FSM struct {
	cm *cache.CacheManager
}

func newFSM() *FSM {
	fm := &FSM{}
	fm.cm = cache.NewCacheManager()
	return fm
}

// Apply applies a Raft log entry to the key-value store.
func (f *FSM) Apply(logEntry *raft.Log) interface{} {
	e := cache.Entry{
	}
	if err := json.Unmarshal(logEntry.Data, &e); err != nil {
		//panic("Failed unmarshaling Raft log entry. This is a bug.")
		return err
	}
	ret := f.cm.Set(e.Menu, e.Key, e.Value)
	fmt.Println("apply", logEntry.Data)
	logrus.Printf("fms.Apply(), logEntry:%s, ret:%v\n", logEntry.Data, ret)
	return ret
}

// Snapshot returns a latest snapshot
func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return &snapshot{cm: f.cm}, nil
}

// Restore stores the key-value store to a previous state.
func (f *FSM) Restore(serialized io.ReadCloser) error {
	return f.cm.UnMarshal(serialized)
}

type snapshot struct {
	cm *cache.CacheManager
}

// Persist saves the FSM snapshot out to the given sink.
func (s *snapshot) Persist(sink raft.SnapshotSink) error {
	snapshotBytes, err := s.cm.Marshal()
	if err != nil {
		sink.Cancel()
		return err
	}

	if _, err := sink.Write(snapshotBytes); err != nil {
		sink.Cancel()
		return err
	}

	if err := sink.Close(); err != nil {
		sink.Cancel()
		return err
	}
	return nil
}

func (s *snapshot) Release() {}
