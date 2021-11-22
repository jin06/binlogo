package factory

type Role uint8

var (
	RoleMaster Role = 01
	RoleWorker Role = 02
	RoleAPI    Role = 04
)

type Factory struct {
	Role Role
}

func (f *Factory) isRole(role Role) bool {
	return f.Role&role == role
}

func (f *Factory) IsMaster() bool {
	return f.isRole(RoleMaster)
}

func (f *Factory) IsAPI() bool {
	return f.isRole(RoleAPI)
}
