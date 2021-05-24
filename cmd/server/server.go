package main

import (
	"flag"
	"fmt"
	"github.com/jin06/binlogo/config"
)

func main() {
	var configPath string
	defaultPath := "./../../config/binlogo.yaml"
	flag.StringVar(&configPath, "config", defaultPath, "config path")
	config.InitCfg(configPath)
	fmt.Println(config.Cfg)
}
