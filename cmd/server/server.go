package main

import (
	"flag"
	"fmt"
	"github.com/jin06/binlogo/config"
	"github.com/jin06/binlogo/store"
	"github.com/sirupsen/logrus"
)

func main() {
	var configPath string
	logrus.SetLevel(logrus.DebugLevel)
	defaultPath := "./../../config/binlogo.yaml"
	flag.StringVar(&configPath, "config", defaultPath, "config path")
	config.InitCfg(configPath)
	store.InitDefault()
	fmt.Println(config.Cfg)
}
