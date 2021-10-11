package app

import (
	pipeline2 "github.com/jin06/binlogo/app/pipeline/pipeline"
	store2 "github.com/jin06/binlogo/pkg/store"
	model2 "github.com/jin06/binlogo/pkg/store/model"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/sirupsen/logrus"
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

	sPipeline := &pipeline.Pipeline{
		Name:      "test",
		AliasName: "本地测试",
		Mysql: &model2.Mysql{
			Address:  "127.0.0.1",
			Port:     13306,
			User:     "root",
			Password: "123456",
			Flavor:   "mysql",
			ServerId: 1001,
		},
		Filters: []*model2.Filter{
			{
			},
		},
		Output: &model2.Output{
			Sender: &model2.Sender{
				Type: model2.SNEDER_TYPE_STDOUT,
				Kafka: &model2.Kafka{
					Brokers: []string{
						"kafka-banana.shan.svc.cluster.local:9092",
					},
					Topic: "binlogo-test1",
				},
			},
		},
	}
	store2.Create(sPipeline)
	//os.Exit(1)

	sPosition := &model2.Position{
		BinlogFile:     "mysql-bin.000014",
		BinlogPosition: 120,
	}
	store2.Create(sPosition)
	p, err := pipeline2.New(
		pipeline2.OptionPipeline(sPipeline),
		pipeline2.OptionPosition(sPosition),
	)
	err = p.Run()

	if err != nil {
		panic(err)
	}
	logrus.Debug("start pipeline")
	select {}
}

