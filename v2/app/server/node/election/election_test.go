package election

import (
	"context"
	"testing"
	"time"

	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/pkg/node/role"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
)

func TestSetRole(t *testing.T) {
	e := New()

	for i := 0; i < 10; i++ {
		rr := role.LEADER
		if i%2 == 0 {
			rr = role.LEADER
		}
		go func() {
			e.SetRole(rr)
		}()
	}

	go func() {
		for i := 0; i < 10; i++ {
			rr := role.LEADER
			if i%2 == 0 {
				rr = role.LEADER
			}
			rrr := <-e.RoleCh
			if rrr != rr {
				t.Fail()
			}
		}
	}()
}

func TestRun(t *testing.T) {
	configs.InitGoTest()
	n := New(OptionNode(node.NewNode("go_test_node")), OptionTTL(5))
	n.Run(context.Background())
	time.Sleep(time.Millisecond * 200)
	l, err := n.Leader()
	if err != nil {
		t.Error(err)
	}
	if l != "go_test_node" {
		t.Log(l)
		t.Fail()
	}
	t.Log(n.Role())

	err = n.Resign(context.Background())
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Millisecond * 200)
}
