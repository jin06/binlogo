package config

import (
	"context"
	"fmt"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
)

var Cfg Config

type Config struct {
	Cluster struct {
		Name string `config:"name"`
	} `config:"cluster"`
	Node struct {
		Name string `config:"name"`
	} `config:"node"`
	Api struct {
		Port uint8 `config:"port"`
	} `config:"api"`
	Store struct {
		Type string `config:"type"`
		Etcd struct {
			Endpoints []string `config:"endpoints"`
		} `config:"etcd"`
	} `config:"store"`
}

func InitCfg(path string) {
	loader := confita.NewLoader(file.NewBackend(path))
	err := loader.Load(context.Background(), &Cfg)
	fmt.Println("config", Cfg)
	if err != nil {
		panic(err)
	}
}
