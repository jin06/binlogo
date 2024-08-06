package ip

import "testing"

func TestLocalIp(t *testing.T) {
	localIp, err := LocalIp()
	if err != nil {
		t.Fail()
	}
	if localIp == nil {
		t.Fail()
	} else {
		if localIp.String() == "" {
			t.Fail()
		}
	}
}
