package app

import (
	pipeline2 "github.com/jin06/binlogo/app/pipeline/pipeline"
)

func RunServer() {
	//var configPath string
	//logrus.SetLevel(logrus.DebugLevel)
	//defaultPath := "./../../config/binlogo.yaml"
	//flag.StringVar(&configPath, "config", defaultPath, "config path")
	//store.InitDefault()
	//etcd.DefaultETCD()

	p2, err := pipeline2.NewFromStore("test")
	if err != nil {
		panic(err)
	}
	err = p2.Run()
	if err != nil {
		panic(err)
	}
	select {}
}

