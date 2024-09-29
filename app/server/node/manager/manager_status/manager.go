package manager_status

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/jin06/binlogo/v2/pkg/store/dao"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
	"github.com/jin06/binlogo/v2/pkg/util/ip"
	"github.com/sirupsen/logrus"
)

// Manager report node status and resource usage
type Manager struct {
	Node      *node.Node
	syncMutex sync.Mutex
	nodeMutex sync.Mutex
}

// NewManager returns a new Manager
func NewManager(n *node.Node) *Manager {
	sm := &Manager{
		Node:      n,
		syncMutex: sync.Mutex{},
		nodeMutex: sync.Mutex{},
	}
	return sm
}

// Run start working
func (m *Manager) Run(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln("manager status panic, ", r)
		}
	}()
	var err error
	statusTicker := time.NewTicker(time.Second * 10)
	defer statusTicker.Stop()
	ipTicker := time.NewTicker(time.Minute)
	defer ipTicker.Stop()
	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		case <-statusTicker.C:
			{
				if err = m.syncStatus(); err != nil {
					logrus.Error("Sync status failed: ", err)
				}
			}
		case <-ipTicker.C:
			{
				if err = m.syncIP(); err != nil {
					logrus.Error("Sync node ip failed: ", err)
				}
			}
		}
	}
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
		if ok, err = dao.UpdateNode(m.Node.Name, node.WithNodeIP(nip)); err != nil {
			return
		}
		if !ok {
			return
		}
		m.Node.IP = nip
	}
	return
}
