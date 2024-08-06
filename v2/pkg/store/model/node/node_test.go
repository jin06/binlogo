package node

import (
	"net"
	"testing"

	"github.com/jin06/binlogo/configs"
)

func TestNode(t *testing.T) {
	node := NewNode("go_test_node")
	if node == nil {
		t.Fail()
	}
	if node.Key() == "" {
		t.Fail()
	}
	node2 := NewNode("go_test_node2")
	err := node2.Unmarshal([]byte(node.Val()))
	if err != nil {
		t.Error(err)
	}
	if node.Name != node2.Name {
		t.Fail()
	}
	localIP := net.IPv4(127, 0, 0, 1)
	WithNodeIP(localIP)(node)
	WithNodeIP(localIP)(node2)
	if !node.IP.Equal(node2.IP) {
		t.Fail()
	}
	WithNodeVersion(configs.VERSITON)(node)
	WithNodeVersion(configs.VERSITON)(node2)
	if node.Version != node2.Version {
		t.Fail()
	}
}
