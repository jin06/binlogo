package app

import (
	"fmt"
	"github.com/jin06/binlogo/config"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/store/etcd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCommand() (cmd *cobra.Command) {
	cmdVersion := &cobra.Command{
		Use:   "version",
		Short: "version",
		Long:  "version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("binlogo version: " + config.VERSITON)
		},
	}

	cmdRun := &cobra.Command{
		Use:   "server",
		Short: "Generate mysql data increment",
		Long:  "Generate mysql data increment",
		Run: func(cmd *cobra.Command, args []string) {
			blog.Env(config.Env(viper.GetString("env")))
			config.InitViperFromFile(viper.GetString("config"))
			etcd.DefaultETCD()
			//RunServer()
			RunNode()
		},
	}

	cmd = &cobra.Command{Use: "binlogo"}
	var cfgFile string
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "./config/binlogo.yaml", "config file default is ./config/binlogo.yaml")
	err := viper.BindPFlag("config", cmd.PersistentFlags().Lookup("config"))
	if err != nil {
		panic(err)
	}
	cmd.AddCommand(cmdRun, cmdVersion)
	return
}
