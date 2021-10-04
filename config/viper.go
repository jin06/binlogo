package config

import (
	"github.com/spf13/viper"
)

//func init() {
//	InitViperFromFile()
//}

func InitViperFromFile() {
	configName := "binlogo"
	configPath := "./config/"
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
