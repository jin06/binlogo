package app

import (
	"context"
	"errors"
	"time"

	node2 "github.com/jin06/binlogo/app/server/node"
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/store/dao/dao_node"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/sirupsen/logrus"
)

// RunNode run node.
func RunNode(c context.Context) (resCtx context.Context, err error) {
	logrus.Info("init node")
	resCtx, cancel := context.WithCancel(c)
	defer func() {
		if err != nil {
			cancel()
		}
	}()
	nModel := &node.Node{
		Name:       configs.NodeName,
		Version:    configs.Version,
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
	go func() {
		logrus.Info("run node")
		cancel()
		for {
			_node := node2.New(node2.OptionNode(nModel))
			ctx, _ := context.WithCancel(c)
			ch := make(chan error, 1)
			go func() {
				var nodeError error
				defer func() {
					if r := recover(); r != nil {
						logrus.Errorln("run node panic, ", r)
						nodeError = errors.New("panic")
					}
					ch <- nodeError
				}()
				logrus.WithField("node name", configs.NodeName).Info("Running node")
				nodeError = _node.Run(ctx)
			}()
			errNode := <-ch
			if errNode != nil {
				logrus.Errorln("run node error, ", errNode)
			}
			close(ch)
			time.Sleep(time.Second * 5)
		}
	}()
	return
}
