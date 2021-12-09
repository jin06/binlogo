package app

import (
	"fmt"

	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/blog"
	store2 "github.com/jin06/binlogo/pkg/store"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewCommand return a new *cobra.Command for cli
func NewCommand() (cmd *cobra.Command) {
	cmd = &cobra.Command{Use: "binctl", Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("Operate pipeline")
		//fmt.Println("init binctl")
	}}
	cmd.AddCommand(cmdMemory())
	cmd.AddCommand(cmdPipeline())
	cmd.PersistentFlags().String("config", "./configs/binlogo.yaml", "config file default is ./configs/binlogo.yaml")
	err := viper.BindPFlag("config", cmd.PersistentFlags().Lookup("config"))
	// configs.InitViperFromFile(viper.GetString("config"))
	configs.Init(viper.GetString("config"))
	//etcd2.DefaultETCD()
	blog.Env(configs.Env(viper.GetString("env")))
	if err != nil {
		fmt.Println(err)
	}
	return
}

func cmdMemory() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "memory",
		Short: "Show server memory usage",
		Long:  "Show server memory usage",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	return
}

func cmdCreatePipe() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "create",
		Short: "Create pipeline",
		Run: func(cmd *cobra.Command, args []string) {
			sPipeline := &pipeline.Pipeline{
				Name:      "test",
				AliasName: "本地测试",
				Mysql: &pipeline.Mysql{
					Address:  "127.0.0.1",
					Port:     13306,
					User:     "root",
					Password: "123456",
					Flavor:   pipeline.FLAVOR_MYSQL,
					ServerId: 1001,
				},
				Filters: []*pipeline.Filter{
					{},
				},
				Output: &pipeline.Output{
					Sender: &pipeline.Sender{
						Type: pipeline.SNEDER_TYPE_STDOUT,
						Kafka: &pipeline.Kafka{
							Brokers: "kafka-banana.shan.svc.cluster.local:9092",
							Topic:   "binlogo-test1",
						},
					},
				},
			}
			ok, err := store2.Create(sPipeline)
			fmt.Println(ok, err)
		},
	}
	return
}

func cmdPipeline() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "pipe",
		Short: "Operate pipeline",
		Long:  "Operate pipeline",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	cmd.AddCommand(cmdCreatePipe())
	return
}
