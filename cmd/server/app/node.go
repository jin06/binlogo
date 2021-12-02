package app

import (
	"context"
	node2 "github.com/jin06/binlogo/app/server/node"
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/store/dao/dao_node"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/sirupsen/logrus"
	"time"
)

// RunNode run node.
func RunNode() (err error) {
	nModel := &node.Node{
		Name:       configs.NodeName,
		Version:    configs.VERSITON,
		CreateTime: time.Now(),
		Role:       node.Role{Master: true, Admin: true, Worker: true},
	}
	nModel.IP = configs.NodeIP
	n, err := dao_node.GetNode(nModel.Name)
	if err != nil {
		return
	}
	if n != nil {
		if _, err = dao_node.UpdateNode(nModel.Name, node.WithNodeIP(nModel.IP), node.WithNodeVersion(nModel.Version)); err != nil {
			return
		}
	} else {
		if err = dao_node.CreateNodeIfNotExist(nModel); err != nil {
			return
		}
	}
	err = dao_node.CreateStatusIfNotExist(&node.Status{NodeName: nModel.Name, Ready: true})
	if err != nil {
		logrus.Error(err)
		return
	}
	var _node *node2.Node
	if _node, err = node2.New(node2.OptionNode(nModel)); err != nil {
		return
	}
	ctx := context.TODO()
	go func() {
		for {
			logrus.WithField("node name", configs.NodeName).Info("Running node")
			err = _node.Run(ctx)
			time.Sleep(time.Second * 5)
		}
	}()

	return
}
