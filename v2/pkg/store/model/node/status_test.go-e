package node

import "testing"

func TestStatus(t *testing.T) {
	s := New("go_test_node")
	WithReady(false)(s)
	if s.Ready == true {
		t.Fail()
	}
	WithNetworkUnavailable(true)(s)
	if s.NetworkUnavailable == false {
		t.Fail()
	}
	WithMemoryPressure(true)(s)
	if s.MemoryPressure == false {
		t.Fail()
	}
	WithDiskPressure(true)(s)
	if s.DiskPressure == false {
		t.Fail()
	}
	WithCPUPressure(true)(s)
	if s.CPUPressure == false {
		t.Fail()
	}
}
