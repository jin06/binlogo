package app

import (
	"context"
	node2 "github.com/jin06/binlogo/app/server/node"
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/store/dao"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/util/ip"
	"github.com/spf13/viper"
)

func RunNode() (err error) {
	nModel := &node.Node{
		Name :viper.GetString("node.name"),
		Status: node.STATUS_ON,
		Version: configs.VERSITON,
	}
	if nModel.Ip == nil {
		if nModel.Ip, err = ip.LocalIp(); err != nil {
			return
		}
	}
	if err = dao.CreateNodeIfNotExist(nModel); err != nil {
		return
	}
	var _node *node2.Node
	if _node, err = node2.New(node2.OptionNode(nModel)); err != nil {
		return
	}
	ctx := context.Background()
	err = _node.Run(ctx)
	return
}
