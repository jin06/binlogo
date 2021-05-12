package config

import (
	"context"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
)

var Cfg Config
var DefaultPath string = "/Users/jinlong/code/jin06/binlogo/config/binlogo.yaml"

type Config struct {
	Cluster struct{
		Name string `config:"name"`
	}
	Node struct{
		Name string `config:"name"`
	}
}

func InitCfg() {
	loader := confita.NewLoader(file.NewBackend(DefaultPath))
	loader.Load(context.Background(), &Cfg)
}