package factory

import "testing"

func TestRole(t *testing.T) {
	f := Factory{}
	f.Role = RoleMaster + RoleWorker
	t.Logf("role: %b ", f.Role)
	if !f.IsMaster() {
		t.Fatal("Expect master")
	}
	if f.IsAPI() {
		t.Fatal("Expect not api")
	}
}
