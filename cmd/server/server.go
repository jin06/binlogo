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
	//defaultPath := "./../../config/binlogo.yaml"
	defaultPath := "config/binlogo.yaml"
	flag.StringVar(&configPath, "config", defaultPath, "config path")
	config.InitCfg(configPath)
	//store.InitDefault()
	etcd.DefaultETCD()
	store.Create(&model.Pipeline{
		ID:   "123",
		Name: "测试",
	})
	p := model.Pipeline{
		ID: "123",
	}
	store.Create(&model.Filter{
		ID: "999",
		PipelineId: "123",
	})
	store.Get(&p)
	fmt.Println(p)
}
