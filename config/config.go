package config

var Cfg Config

type Config struct {
	Env     Env `config:"env"`
	Cluster struct {
		Name string `config:"name"`
	} `config:"cluster"`
	Node struct {
		Name string `config:"name"`
	} `config:"node"`
	Api struct {
		Port uint `config:"port"`
	} `config:"api"`
	Store struct {
		Type string `config:"type"`
		Etcd struct {
			Endpoints []string `config:"endpoints"`
		} `config:"etcd"`
	} `config:"store"`
}

type Env string

const (
	ENV_PRO  = "production"
	ENV_DEV  = "dev"
	ENV_TEST = "test"
)
