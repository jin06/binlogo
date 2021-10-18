package app

import (
	"context"
	"fmt"
	node2 "github.com/jin06/binlogo/app/server/node"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/spf13/viper"
	"os"
)

func RunNode() {
	nModel := &node.Node{
	}
	nModel.Name = viper.GetString("node.name")
	_node, err := node2.New(node2.OptionNode(nModel))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	ctx := context.Background()
	err = _node.Run(ctx)
	if err != nil {
		fmt.Println(err)
	}
}
