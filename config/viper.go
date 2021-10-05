package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

//func init() {
//	InitViperFromFile()
//}

func InitViperFromFile(file string) {
	if file != "" {
		viper.SetConfigFile(file)
	} else {
		file = "./config/binlogo.yaml"
	}
	if err := viper.ReadInConfig(); err != nil {
		logrus.Error(err.Error())
		os.Exit(1)
	}
}
