package configs

import (
	"github.com/jin06/binlogo/pkg/util/ip"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"syscall"
)

// Init init configs
func Init(file string) {
	initViperDefault()
	initViperFromFile(file)
	initViperFromEnv()
	initConst()
}

// initViperFromFile read config file and write to viper
func initViperFromFile(file string) {
	if file == "" {
		file = "./configs/binlogo.yaml"
	}
	//fmt.Println("init config from ", file)
	viper.SetConfigFile(file)
	if err := viper.ReadInConfig(); err != nil {
		logrus.Error(err.Error())
		//os.Exit(1)
	}
}

// initViperDefault set default configs
func initViperDefault() {
	hostName, _ := os.Hostname()
	viper.SetDefault("node.name", hostName)
	viper.SetDefault("env", ENV_PRO)
	viper.SetDefault("cluster.name", CLUSTER_NAME)
	viper.SetDefault("console.listen", CONSOLE_LISTEN)
	viper.SetDefault("console.port", CONSOLE_PORT)
}

// initViperFromEnv read config from env then whrite to viper
func initViperFromEnv() {
	if val, found := syscall.Getenv("NODE_NAME"); found {
		viper.Set("node.name", val)
	}
	if val, found := syscall.Getenv("BINLOGO_ENV"); found {
		viper.Set("env", val)
	}
	if val, found := syscall.Getenv("CLUSTER_NAME"); found {
		viper.Set("cluster.name", val)
	}
	if val, found := syscall.Getenv("CONSOLE_LISTEN"); found {
		viper.Set("console.listen", val)
	}
	if val, found := syscall.Getenv("CONSOLE_PORT"); found {
		viper.Set("console.port", val)
	}
	if val, found := syscall.Getenv("ETCD_ENDPOINTS"); found {
		viper.Set("etcd.endpoints", val)
	}
	if val, found := syscall.Getenv("ETCD_PASSWORD"); found {
		viper.Set("etcd.password", val)
	}
	if val, found := syscall.Getenv("ETCD_USERNAME"); found {
		viper.Set("etcd.username", val)
	}
}

// initConst set global config
func initConst() {
	ENV = Env(viper.GetString("env"))
	NodeName = viper.GetString("node.name")
	NodeIP, _ = ip.LocalIp()
}

// InitGoTest init environment for testing
func InitGoTest() {
	_ = os.Setenv("NODE_NAME", "go_test_node")
	_ = os.Setenv("BINLOGO_ENV", "dev")
	_ = os.Setenv("CLUSTER_NAME", "go_test_cluster")
	_ = os.Setenv("CONSOLE_LISTEN", "0.0.0.0")
	_ = os.Setenv("CONSOLE_PORT", "19999")
	_ = os.Setenv("ETCD_ENDPOINTS", "localhost:12379")
	_ = os.Setenv("ETCD_PASSWORD", "")
	_ = os.Setenv("ETCD_USERNAME", "")
	Init("")
}
