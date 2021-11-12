package app

import (
	"fmt"
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/blog"
	etcd2 "github.com/jin06/binlogo/pkg/store/etcd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
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
			cfg, _ := cmd.Flags().GetString("config")
			configs.InitViperFromFile(cfg)
			etcd2.DefaultETCD()
			configs.ENV = configs.Env(viper.GetString("env"))
			//blog.Env(configs.Env(viper.GetString("env")))
			blog.Env(configs.ENV)
			//RunServer()
			logrus.Infoln("init configs finish")
			if configs.ENV == configs.ENV_DEV {
				go func() {
					logrus.Println(http.ListenAndServe("localhost:6060",nil))
				}()
			}
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
	//var cfgFile string
	//cmd.PersistentFlags().StringVar(&cfgFile, "config", "./configs/binlogo.yaml", "config file default is ./configs/binlogo.yaml")
	cmd.PersistentFlags().String("config","./configs/binlogo.yaml","")
	//panic(123)

	//if err != nil {
	//	panic(err)
	//}
	cmd.AddCommand(cmdRun, cmdVersion)
	return
}
