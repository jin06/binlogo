package node

type Role struct {
	Master bool `json:"master"`
	Admin  bool `json:"admin"`
	Worker bool `json:"worker"`
}
