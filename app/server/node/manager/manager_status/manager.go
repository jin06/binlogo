package manager_status

import (
	"context"
	"github.com/jin06/binlogo/pkg/store/dao/dao_node"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/util/ip"
	"github.com/sirupsen/logrus"
	"net"
	"sync"
	"time"
)

type Manager struct {
	Node      *node.Node
	syncMutex sync.Mutex
	nodeMutex sync.Mutex
}

func NewManager(n *node.Node) *Manager {
	sm := &Manager{
		Node:      n,
		syncMutex: sync.Mutex{},
		nodeMutex: sync.Mutex{},
	}
	return sm
}

func (m *Manager) Run(ctx context.Context) (err error) {
	if err = m.syncStatus(); err != nil {
		logrus.Error("Sync status failed: ", err)
		return
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			case <-time.Tick(time.Second * 10):
				//case <-time.Tick(time.Second * 1):
				{
					if err = m.syncStatus(); err != nil {
						logrus.Error("Sync status failed: ", err)
					}
				}
			case <-time.Tick(time.Minute):
				{
					if err = m.syncIP(); err != nil {
						logrus.Error("Sync node ip failed: ", err)
					}
				}
			}
		}
	}()
	return
}

func (m *Manager) syncStatus() (err error) {
	//m.syncMutex.Lock()
	//defer m.syncMutex.Unlock()
	ns := NewNodeStatus(m.Node.Name)
	err = ns.syncNodeStatus()
	return
}

func (m *Manager) syncIP() (err error) {
	//m.nodeMutex.Lock()
	//defer m.nodeMutex.Unlock()
	var nip net.IP
	if nip, err = ip.LocalIp(); err != nil {
		return
	}
	if m.Node.IP.String() != nip.String() {
		var ok bool
		if ok, err = dao_node.UpdateNode(m.Node.Name, node.WithNodeIP(nip)); err != nil {
			return
		}
		if !ok {
			return
		}
		m.Node.IP = nip
	}
	return
}
