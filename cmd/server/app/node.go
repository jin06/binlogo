package app

import (
	"context"
	"fmt"
	"time"

	server_node "github.com/jin06/binlogo/v2/app/server/node"
	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/pkg/event"
	"github.com/jin06/binlogo/v2/pkg/store/dao"
	model_event "github.com/jin06/binlogo/v2/pkg/store/model/event"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
	"github.com/sirupsen/logrus"
)

// RunNode run node.
func RunNode(c context.Context) (err error) {
	logrus.Info("init node")
	nodeOpts := &node.Node{
		Name:       configs.Default.NodeName,
		Version:    configs.Version,
		CreateTime: time.Now(),
		Role:       node.Role{Master: true, Admin: true, Worker: true},
	}
	nodeOpts.IP = configs.NodeIP
	n, err := dao.GetNode(c, nodeOpts.Name)
	if err != nil {
		return
	}
	if n != nil {
		if _, err = dao.UpdateNode(c, nodeOpts.Name, node.WithNodeIP(nodeOpts.IP), node.WithNodeVersion(nodeOpts.Version)); err != nil {
			return
		}
	} else {
		if err = dao.CreateNodeIfNotExist(c, nodeOpts); err != nil {
			return
		}
	}
	err = dao.CreateStatusIfNotExist(&node.Status{NodeName: nodeOpts.Name, Ready: true})
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Info("run node")
	event.Event(model_event.NewInfoNode("Run node success"))
	for {
		server := server_node.New(server_node.OptionNode(nodeOpts))
		func() {
			var err error
			ctx, cancel := context.WithCancel(c)
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("server panic %v", r)
				}
				if err != nil {
					logrus.Errorln(err)
				}
				cancel()
			}()
			logrus.WithField("node name", configs.NodeName).Info("Running node")
			err = server.Run(ctx)
		}()
		time.Sleep(time.Second * 5)
	}
}
