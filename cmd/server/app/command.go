package app

import (
	"fmt"
	"os"

	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/event"
	event2 "github.com/jin06/binlogo/pkg/store/model/event"
	"github.com/spf13/cobra"
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
			Init(cfg)
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
