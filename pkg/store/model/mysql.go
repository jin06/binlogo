package model

type Mysql struct {
	Address    string `json:"address"`
	Port       uint16 `json:"post"`
	User       string `json:"user"`
	Password   string `json:"password"`
	ServerId   uint32 `json:"server_id"`
	Flavor     string `json:"flavor"`
}
