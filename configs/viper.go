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

}

func InitConfigs() {
	node_name := os.Getenv("NODE_NAME")
	if node_name == "" {
		node_name, _ = os.Hostname()
	}
	viper.SetDefault("node.name", node_name)
	viper.SetDefault("env", os.Getenv("BINLOGO_ENV"))
	viper.SetDefault("cluster.name", os.Getenv("CLUSTER_NAME"))
	consoleListen := os.Getenv("CONSOLE_LISTEN")
	if consoleListen == "" {
		consoleListen = CONSOLE_LISTEN
	}
	viper.SetDefault("console.listen", consoleListen)
	consolePort := os.Getenv("CONSOLE_PORT")
	if consolePort == "" {
		consolePort = CONSOLE_PORT
	}
	viper.SetDefault("etcd.endpoints", os.Getenv("ETCD_ENDPOINTS"))
	viper.SetDefault("etcd.password", os.Getenv("ETCD_PASSWORD"))
	viper.SetDefault("console.port", consolePort)

	ENV = Env(viper.GetString("env"))
	NodeName = viper.GetString("node.name")
	NodeIP, _ = ip.LocalIp()
}
