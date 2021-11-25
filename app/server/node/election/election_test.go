package election

import (
	"github.com/jin06/binlogo/pkg/node/role"
	"testing"
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
