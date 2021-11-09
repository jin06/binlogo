package app

import (
	"fmt"
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/blog"
	etcd2 "github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func NewCommand() (cmd *cobra.Command) {
	cmdVersion := &cobra.Command{
		Use:   "version",
		Short: "version",
		Long:  "version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("binlogo version: " + configs.VERSITON)
		},
	}

	cmdRun := &cobra.Command{
		Use:   "server",
		Short: "Generate mysql data increment",
		Long:  "Generate mysql data increment",
		Run: func(cmd *cobra.Command, args []string) {
			configs.InitViperFromFile(viper.GetString("config"))
			etcd2.DefaultETCD()
			blog.Env(configs.Env(viper.GetString("env")))
			//RunServer()
			logrus.Infoln("init configs finish")
			if err := RunNode(); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			if err := RunConsole(); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			select {}
		},
	}

	cmd = &cobra.Command{Use: "binlogo"}
	var cfgFile string
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "./configs/binlogo.yaml", "config file default is ./config/binlogo.yaml")
	err := viper.BindPFlag("config", cmd.PersistentFlags().Lookup("config"))
	if err != nil {
		panic(err)
	}
	cmd.AddCommand(cmdRun, cmdVersion)
	return
}
