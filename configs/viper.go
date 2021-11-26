package configs

import (
	"github.com/jin06/binlogo/pkg/util/ip"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

// InitViperFromFile read config file and write to viper
func InitViperFromFile(file string) {
	if file == "" {
		file = "./configs/binlogo.yaml"
	}
	//fmt.Println("init config from ", file)
	viper.SetConfigFile(file)
	if err := viper.ReadInConfig(); err != nil {
		logrus.Error(err.Error())
		os.Exit(1)
	}
}

// InitConfigs sets some default values to viper from system environment
func InitConfigs() {
	nodeName := os.Getenv("NODE_NAME")
	if nodeName == "" {
		nodeName, _ = os.Hostname()
	}
	viper.SetDefault("node.name", nodeName)
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
	viper.SetDefault("console.port", consolePort)
	viper.SetDefault("etcd.endpoints", os.Getenv("ETCD_ENDPOINTS"))
	viper.SetDefault("etcd.password", os.Getenv("ETCD_PASSWORD"))
	ENV = Env(viper.GetString("env"))
	NodeName = viper.GetString("node.name")
	NodeIP, _ = ip.LocalIp()
}

// InitEnv sets some default values to environment for testing
func DefaultEnv() {
	_ = os.Setenv("NODE_NAME", "go_test_node")
	_ = os.Setenv("BINLOGO_ENV", "dev")
	_ = os.Setenv("CLUSTER_NAME", "go_test_cluster")
	_ = os.Setenv("CONSOLE_LISTEN", "0.0.0.0")
	_ = os.Setenv("CONSOLE_PORT", "19999")
	_ = os.Setenv("ETCD_ENDPOINTS", "localhost:12379")
	_ = os.Setenv("ETCD_PASSWORD", "")
}

// InitGoTest init environment for testing
func InitGoTest() {
	DefaultEnv()
	InitConfigs()
}

