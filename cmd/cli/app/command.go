package app

import (
	"github.com/jin06/binlogo/pkg/ps"
	"github.com/spf13/cobra"
)

func NewCommand() (cmd *cobra.Command) {
	cmd = &cobra.Command{Use:"app"}
	cmd.AddCommand(cmdMemory())
	return
}

func cmdMemory() (cmd *cobra.Command) {
	cmd =  &cobra.Command{
		Use: "memory",
		Short: "Show server memory usage",
		Long:"Show server memory usage",
		Run: func(cmd *cobra.Command, args []string) {
			ps.Memory()
		},
	}
	return
}
