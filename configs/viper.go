package configs

import (
	"fmt"
	"github.com/jin06/binlogo/pkg/util/ip"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

//func init() {
//	InitViperFromFile()
//}

func InitViperFromFile(file string) {
	if file == "" {
		file = "./configs/binlogo.yaml"
	}
	fmt.Println("init config from ", file)
	viper.SetConfigFile(file)
	if err := viper.ReadInConfig(); err != nil {
		logrus.Error(err.Error())
		os.Exit(1)
	}
	viper.SetDefault("console.listen", CONSOLE_LISTEN)
	viper.SetDefault("console.port", CONSOLE_PORT)
}

func InitConfigs() {
	hostname, _ := os.Hostname()
	viper.SetDefault("node.name", hostname)
	ENV = Env(viper.GetString("env"))
	NodeName = viper.GetString("node.name")
	NodeIP, _  = ip.LocalIp()
}
