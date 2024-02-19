package app

import (
	"context"
	"fmt"
	"sync"

	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/promeths"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		Run: func(cmd *cobra.Command, args []string) {
			cfg, _ := cmd.Flags().GetString("config")
			Init(cfg)
			promeths.Init()
			RunEvent()
			// panic(123)
			var err error
			ctx := context.Background()
			exit := make(chan struct{})
			closeOnce := sync.Once{}
			go func() {
				if err = RunNode(ctx); err != nil {
					fmt.Println(err.Error())
				}
				closeOnce.Do(func() {
					close(exit)
				})
			}()
			if viper.GetBool("roles.api") {
				go func() {
					if err = RunConsole(ctx); err != nil {
						fmt.Println(err.Error())
					}
					closeOnce.Do(func() {
						close(exit)
					})
				}()
			}
			<-exit
		},
	}

	cmd = &cobra.Command{Use: "binlogo"}
	cmd.PersistentFlags().String("config", "./etc/binlogo.yaml", "")
	cmd.AddCommand(cmdRun, cmdVersion)
	return
}
