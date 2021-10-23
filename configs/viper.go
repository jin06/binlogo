package configs

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
		file = "./configs/binlogo.yaml"
	}
	if err := viper.ReadInConfig(); err != nil {
		logrus.Error(err.Error())
		os.Exit(1)
	}
	viper.SetDefault("console.listen" , CONSOLE_LISTEN)
	viper.SetDefault("console.port", CONSOLE_PORT)
}
