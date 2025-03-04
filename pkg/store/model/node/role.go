package node

// Role store struct
type Role struct {
	Master bool `json:"master" redis:"master"`
	Admin  bool `json:"admin" redis:"admin"`
	Worker bool `json:"worker" redis:"worker"`
}
