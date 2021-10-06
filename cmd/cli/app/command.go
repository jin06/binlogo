package app

import (
	"fmt"
	"github.com/jin06/binlogo/pkg/ps"
	"github.com/spf13/cobra"
)

func NewCommand() (cmd *cobra.Command) {
	cmd = &cobra.Command{Use:"binctl"}
	cmd.AddCommand(cmdMemory())
	cmd.AddCommand(cmdPipeline())
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


func cmdPipeline() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use: "pipe",
		Short: "Operate pipeline",
		Long:"Operate pipeline",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Operate pipeline")
			str, _ := cmd.PersistentFlags().GetString("f")
			fmt.Println(str)
		},
	}
	cmd.PersistentFlags().String("f", "","Apply pipeline")
	return
}
