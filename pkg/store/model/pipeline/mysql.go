package pipeline

import "github.com/go-mysql-org/go-mysql/mysql"

type Mysql struct {
	Address  string `json:"address"`
	Port     uint16 `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	ServerId uint32 `json:"server_id"`
	Flavor   Flavor `json:"flavor"`
	Mode     Mode   `json:"mode"`
}
type Mode string

const (
	MODE_GTID    Mode = "gtid"
	MODE_POISTION Mode = "position"
)

type Flavor string

const (
	FLAVOR_MYSQL   Flavor = "MySQL"
	FLAVOR_MARIADB Flavor = "MariaDB"
)

func (f Flavor) YaString() string {
	if f == FLAVOR_MARIADB {
		return mysql.MariaDBFlavor
	}
	return mysql.MySQLFlavor
}
