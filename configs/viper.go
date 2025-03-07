package configs

import (
	"os"
	"strconv"
	"syscall"

	"github.com/jin06/binlogo/v2/pkg/util/ip"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Init init configs
func Init(file string) error {
	if err := InitConfig(file); err != nil {
		return err
	}
	if err := initConfigFromEnv(); err != nil {
		return err
	}
	initViperDefault()
	initViperFromFile(file)
	initConst()
	return nil
}

func initConfigFromEnv() error {
	if val, found := syscall.Getenv("CLUSTER_NAME"); found {
		Default.ClusterName = val
	}
	if val, found := syscall.Getenv("NODE_NAME"); found {
		Default.NodeName = val
	}
	if val, found := syscall.Getenv("REDIS_ADDR"); found {
		Default.Store.Redis.Addr = val
	}
	if val, found := syscall.Getenv("REDIS_PORT"); found {
		if port, err := strconv.Atoi(val); err != nil {
			return err
		} else {
			Default.Store.Redis.Port = port
		}
	}
	if val, found := syscall.Getenv("REDIS_PASSWORD"); found {
		Default.Store.Redis.Password = val
	}
	if val, found := syscall.Getenv("REDIS_DB"); found {
		if port, err := strconv.Atoi(val); err != nil {
			return err
		} else {
			Default.Store.Redis.DB = port
		}
	}
	if val, found := syscall.Getenv("AUTH_TYPE"); found {
		Default.Auth.Type = val
	}
	if val, found := syscall.Getenv("CONSOLE_LISTEN"); found {
		Default.Console.Listen = val
	}
	if val, found := syscall.Getenv("CONSOLE_PORT"); found {
		Default.Console.Port, _ = strconv.Atoi(val)
	}
	if val, found := syscall.Getenv("ROLES_API"); found {
		Default.Roles.API, _ = strconv.ParseBool(val)
	}
	if val, found := syscall.Getenv("ROLES_MASTER"); found {
		Default.Roles.Master, _ = strconv.ParseBool(val)
	}
	return nil
}

// initViperFromFile read config file and write to viper
func initViperFromFile(file string) {
	if file == "" {
		file = "./etc/binlogo.yaml"
	}
	//fmt.Println("init config from ", file)
	logrus.Info("config file path: ", file)
	viper.SetConfigFile(file)
	if err := viper.ReadInConfig(); err != nil {
		logrus.Error(err.Error())
		//os.Exit(1)
	}
}

// initViperDefault set default configs
func initViperDefault() {
	viper.SetDefault("roles.api", true)
	viper.SetDefault("roles.master", true)
	viper.SetDefault("monitor.port", 8085)
}

// initConst set global config
func initConst() {
	NodeIP, _ = ip.LocalIp()
}

// InitGoTest init environment for testing
func InitGoTest() {
	_ = os.Setenv("NODE_NAME", "go_test_node")
	_ = os.Setenv("BINLOGO_ENV", "dev")
	_ = os.Setenv("CLUSTER_NAME", "go_test_cluster")
	_ = os.Setenv("CONSOLE_LISTEN", "0.0.0.0")
	_ = os.Setenv("CONSOLE_PORT", "19999")
	Init("")
}
