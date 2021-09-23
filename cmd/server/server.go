package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jin06/binlogo/config"
	"github.com/jin06/binlogo/pipeline/pipeline"
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
	sPipeline := &model.Pipeline{
		ID:   pipeId,
		Name: "本地测试",
	}
	store.Create(sPipeline)
	sMysql := &model.Mysql{
		Address:    "127.0.0.1",
		Port:       13306,
		User:       "root",
		Password:   "123456",
		Flavor:     "mysql",
		PipelineId: pipeId,
		ServerId:   1001,
	}
	store.Create(sMysql)

	sFilter := &model.Filter{
		ID:         "1",
		PipelineId: pipeId,
	}

	store.Create(sFilter)

	sPosition := &model.Position{
		BinlogFile:     "mysql-bin.000014",
		BinlogPosition: 120,
		PipelineID:     pipeId,
	}
	store.Create(sPosition)
	p, err := pipeline.New(
		pipeline.OptionPipeline(sPipeline),
		pipeline.OptionMysql(sMysql),
		pipeline.OptionPosition(sPosition),
	)
	err = p.Run()

	if err != nil {
		panic(err)
	}
	logrus.Debug("start pipeline")
	go func() {
		for {
			msg := <-p.Input.Ch
			//fmt.Println("println message", msg)
			b, _ := json.Marshal(msg)
			if msg != nil {
				fmt.Println(msg.Content.Head.Type)
			}
			fmt.Println(string(b))
		}
	}()
	select {}
}
