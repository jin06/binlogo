package app

import (
	"context"
	node2 "github.com/jin06/binlogo/app/server/node"
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/store/dao/dao_node"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/util/ip"
	"github.com/spf13/viper"
	"time"
)

func RunNode() (err error) {
	nModel := &node.Node{
		Name:       viper.GetString("node.name"),
		Version:    configs.VERSITON,
		CreateTime: time.Now(),
		Role:       node.Role{Master: true, Admin: true, Worker: true},
	}
	if nModel.IP == nil {
		if nModel.IP, err = ip.LocalIp(); err != nil {
			return
		}
	}
	n, err := dao_node.GetNode(nModel.Name)
	if err != nil {
		return
	}
	if n != nil {
		if _, err = dao_node.UpdateNode(nModel.Name, dao_node.WithNodeIP(nModel.IP), dao_node.WithNodeVersion(nModel.Version)); err != nil {
			return
		}
	} else {
		if err = dao_node.CreateNodeIfNotExist(nModel); err != nil {
			return
		}
	}
	err = dao_node.CreateStatusIfNotExist(nModel.Name, &node.Status{Ready: true})
	if err != nil {
		blog.Error(err)
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
