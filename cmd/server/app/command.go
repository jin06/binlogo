package app

import (
	"context"
	"fmt"

	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/pkg/promeths"
	"github.com/jin06/binlogo/v2/pkg/store/dao"
	store_redis "github.com/jin06/binlogo/v2/pkg/store/redis"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// NewCommand returns a new command for binlogo server
func NewCommand() (cmd *cobra.Command) {
	cmdVersion := &cobra.Command{
		Use:   "version",
		Short: "version",
		Long:  "version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Version: " + configs.Version)
			fmt.Println("Build Time: " + configs.BuildTime)
			fmt.Println("Go Version: " + configs.GoVersion)
		},
	}

	cmdRun := &cobra.Command{
		Use:   "server",
		Short: "Generate mysql data increment",
		Long:  "Generate mysql data increment",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			cfg, _ := cmd.Flags().GetString("config")
			if err := Init(cfg); err != nil {
				return err
			}
			if err := store_redis.Init(ctx, configs.Default.Store.Redis); err != nil {
				return err
			}
			dao.Init()
			promeths.Init()
			RunEvent()
			// panic(123)
			exit := make(chan error, 2)
			go func() {
				if err := RunNode(ctx); err != nil {
					logrus.Error(err)
					panic(err)
					exit <- err
				}
			}()
			if configs.Default.Roles.API {
				go func() {
					if err := RunConsole(ctx); err != nil {
						logrus.Error(err)
						panic(err)
						exit <- err
					}
				}()
			}
			err := <-exit
			return err
		},
	}

	cmd = &cobra.Command{Use: "binlogo"}
	cmd.PersistentFlags().String("config", "./etc/binlogo.yaml", "")
	cmd.AddCommand(cmdRun, cmdVersion)
	return
}
