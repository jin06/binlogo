package dao_cluster

import (
	"testing"

	"github.com/jin06/binlogo/configs"
)

func TestLeaderNode(t *testing.T) {
	configs.InitGoTest()
	name, err := LeaderNode()
	if err != nil {
		t.Error(err)
	}
	t.Log("Leader node: ", name)
}

func TestAllElections(t *testing.T) {
	res, err := AllElections()
	if err != nil {
		t.Error(err)
	}
	t.Log("all elections: ", res)
}
