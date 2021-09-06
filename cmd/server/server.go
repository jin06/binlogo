package main

import (
	"flag"
	"fmt"
	"github.com/jin06/binlogo/config"
	"github.com/jin06/binlogo/store"
	"github.com/jin06/binlogo/store/etcd"
	"github.com/jin06/binlogo/store/model"
	"github.com/sirupsen/logrus"
)

func main() {
	var configPath string
	logrus.SetLevel(logrus.DebugLevel)
	defaultPath := "./../../config/binlogo.yaml"
	flag.StringVar(&configPath, "config", defaultPath, "config path")
	config.InitCfg(configPath)
	//store.InitDefault()
	etcd.DefaultETCD()
	store.Create(&model.Pipeline{
		ID: "123",
		MysqlID: "kkk",
		Database: "sxx_product3",
		Name:"测试",
		Password: "123",
	})
	fmt.Println(config.Cfg)
}
