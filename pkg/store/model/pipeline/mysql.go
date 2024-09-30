package pipeline

import "github.com/go-mysql-org/go-mysql/mysql"

// Mysql store struct
type Mysql struct {
	Address  string `json:"address" redis:"mysql"`
	Port     uint16 `json:"port" redis:"port"`
	User     string `json:"user" redis:"user"`
	Password string `json:"password" redis:"password"`
	ServerId uint32 `json:"server_id" redis:"serverId"`
	Flavor   Flavor `json:"flavor" redis:"flavor"`
	Mode     Mode   `json:"mode" redis:"mode"`
}

// Mode mysql replication mode
type Mode string

const (
	// MODE_GTID GTID Mode
	MODE_GTID Mode = "gtid"
	// MODE_POSITION common Mode
	MODE_POSITION Mode = "position"
)

// Flavor mysql or mariaDB
type Flavor string

const (
	// FLAVOR_MYSQL MySQL DB
	FLAVOR_MYSQL Flavor = "MySQL"
	// FLAVOR_MARIADB MariaDB
	FLAVOR_MARIADB Flavor = "MariaDB"
)

// YaString convert binlogo flavor string
func (f Flavor) YaString() string {
	if f == FLAVOR_MARIADB {
		return mysql.MariaDBFlavor
	}
	return mysql.MySQLFlavor
}
