package app

import (
	server2 "github.com/jin06/binlogo/cmd/app/server"
	"github.com/spf13/cobra"
)

func NewCommand() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "binlogo server",
		Short: "Generate mysql data increment",
		Long:  "Generate mysql data increment",
		Run: func(cmd *cobra.Command, args []string) {
			server2.RunServer()
		},
	}
	return
}
