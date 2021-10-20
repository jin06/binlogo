package status

import (
	"context"
	"github.com/jin06/binlogo/pkg/blog"
	"sync"
	"time"
)

type Manager struct {
	NodeName  string
	syncMutex sync.Mutex
}

func (m *Manager) Run(ctx context.Context) (err error) {
	ns := NewNodeStatus(m.NodeName)
	if err = m.syncStatus(ns); err != nil {
		blog.Error("sync node error: ", err)
		return
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			case <-time.Tick(time.Second * 30):
				{
					if err = m.syncStatus(NewNodeStatus(m.NodeName)); err != nil {
						blog.Error("Sync status failed: ", err)
					}
				}
			}
		}
	}()
	return
}

func (m *Manager) syncStatus(ns *NodeStatus) (err error) {
	m.syncMutex.Lock()
	defer m.syncMutex.Unlock()
	err = ns.syncNodeStatus()
	return
}
