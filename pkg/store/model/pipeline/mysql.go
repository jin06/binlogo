package pipeline

type Mysql struct {
	Address    string `json:"address"`
	Port       uint16 `json:"port"`
	User       string `json:"user"`
	Password   string `json:"password"`
	ServerId   uint32 `json:"server_id"`
	Flavor     string `json:"flavor"`
}
