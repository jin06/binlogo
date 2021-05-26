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
	key := "/binlogo/cluster1/database/1"
	store.Put(key, "klsdjflkdj")
	//fmt.Println(resp)
	fmt.Println(config.Cfg)
}
