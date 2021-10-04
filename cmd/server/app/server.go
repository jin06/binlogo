package app

import (
	"github.com/jin06/binlogo/pipeline/pipeline"
	"github.com/jin06/binlogo/store"
	"github.com/jin06/binlogo/store/model"
	"github.com/sirupsen/logrus"
)

func RunServer() {
	//var configPath string
	//logrus.SetLevel(logrus.DebugLevel)
	//defaultPath := "./../../config/binlogo.yaml"
	//flag.StringVar(&configPath, "config", defaultPath, "config path")
	//store.InitDefault()
	//etcd.DefaultETCD()

	p2, err := pipeline.NewFromStore("test")
	if err != nil {
		panic(err)
	}
	err = p2.Run()
	if err != nil {
		panic(err)
	}
	select {}

	sPipeline := &model.Pipeline{
		Name:      "test",
		AliasName: "本地测试",
		Mysql: &model.Mysql{
			Address:  "127.0.0.1",
			Port:     13306,
			User:     "root",
			Password: "123456",
			Flavor:   "mysql",
			ServerId: 1001,
		},
		Filters: []*model.Filter{
			{
			},
		},
		Output: &model.Output{
			Sender: &model.Sender{
				Type: model.SNEDER_TYPE_STDOUT,
				Kafka: &model.Kafka{
					Brokers: []string{
						"kafka-banana.shan.svc.cluster.local:9092",
					},
					Topic: "binlogo-test1",
				},
			},
		},
	}
	store.Create(sPipeline)
	//os.Exit(1)

	sPosition := &model.Position{
		BinlogFile:     "mysql-bin.000014",
		BinlogPosition: 120,
	}
	store.Create(sPosition)
	p, err := pipeline.New(
		pipeline.OptionPipeline(sPipeline),
		pipeline.OptionPosition(sPosition),
	)
	err = p.Run()

	if err != nil {
		panic(err)
	}
	logrus.Debug("start pipeline")
	select {}
}

