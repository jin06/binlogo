package app

import (
	"fmt"
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/event"
	etcd2 "github.com/jin06/binlogo/pkg/store/etcd"
	event2 "github.com/jin06/binlogo/pkg/store/model/event"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
	"os"
)

// NewCommand returns a new command for binlogo server
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
			configs.InitConfigs()
			//hostname, _ := os.Hostname()
			//viper.SetDefault("node.name", hostname)
			etcd2.DefaultETCD()
			//configs.ENV = configs.Env(viper.GetString("env"))
			blog.Env(configs.ENV)
			logrus.Infoln("init configs finish")
			if configs.ENV == configs.ENV_DEV {
				go func() {
					logrus.Println(http.ListenAndServe("localhost:6060", nil))
				}()
			}
			RunEvent()
			if err := RunNode(); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			event.Event(event2.NewInfoNode("Run node success"))
			if err := RunConsole(); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			select {}
		},
	}

	cmd = &cobra.Command{Use: "binlogo"}
	cmd.PersistentFlags().String("config", "./configs/binlogo.yaml", "")
	cmd.AddCommand(cmdRun, cmdVersion)
	return
}
