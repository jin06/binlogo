package config

import (
	"context"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
)

var Cfg Config

type Config struct {
	Cluster struct{
		Name string `config:"name"`
	}
	Node struct{
		Name string `config:"name"`
	}
	Store struct{
		Type string `config:"type"`
		Etcd struct{
			Endpoints []string `config:"endpoints"`
		}
	}
}

func InitCfg(path string) {
	loader := confita.NewLoader(file.NewBackend(path))
	loader.Load(context.Background(), &Cfg)
}