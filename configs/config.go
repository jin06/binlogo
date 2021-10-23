package configs

import "time"

var Cfg Config

type Config struct {
	Env     Env `config:"env"`
	Cluster struct {
		Name string `config:"name"`
	} `config:"cluster"`
	Node struct {
		Name string `config:"name"`
	} `config:"node"`
	Console struct {
		Port uint `config:"port"`
		Listen string `config:"listen"`
	} `config:"console"`
	Store struct {
		Type string `config:"type"`
		Etcd struct {
			Endpoints []string `config:"endpoints"`
		} `config:"etcd"`
		DialTimeout time.Duration `config:"dial_timeout"`
	} `config:"store"`
}

