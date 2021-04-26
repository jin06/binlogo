package repl

import "github.com/siddontang/go-mysql/replication"

type Config struct {
	ServerID uint32
	Master Master
	Position Position
}

func (c *Config) BinlogSyncerConfig() replication.BinlogSyncerConfig{
	cfg := replication.BinlogSyncerConfig{
		ServerID: c.ServerID,
		Flavor: c.Master.Flavor,
		Host: c.Master.Host,
		Port: c.Master.Port,
		User: c.Master.User,
		Password: c.Master.Password,
	}
	return cfg
}
