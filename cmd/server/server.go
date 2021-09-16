package main

import (
	"flag"
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
	pipeId := "123"
	store.Create(&model.Pipeline{
		ID:   pipeId,
		Name: "本地测试",
	})
	store.Create(&model.Mysql{
		Address:    "127.0.0.1",
		Port:       3306,
		User:       "root",
		Password:   "123456",
		PipelineId: pipeId,
		ServerId:   1001,
	})

	store.Create(&model.Filter{
		ID:         "1",
		PipelineId: pipeId,
	})
	store.Create(&model.Position{
		BinlogFile:     "mysql-bin.000001",
		BinlogPosition: 120,
		PipelineID:     pipeId,
	})
}
