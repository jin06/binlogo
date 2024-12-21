package app

import (
	"context"
	"fmt"
	"time"

	server_node "github.com/jin06/binlogo/v2/app/server/node"
	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/pkg/event"
	"github.com/jin06/binlogo/v2/pkg/store/dao"
	model_event "github.com/jin06/binlogo/v2/pkg/store/model"
	"github.com/jin06/binlogo/v2/pkg/store/model/node"
	"github.com/jin06/binlogo/v2/pkg/util/ip"
	"github.com/sirupsen/logrus"
)

// RunNode run node.
func RunNode(c context.Context) (err error) {
	logrus.Info("init node")
	nip, err := ip.LocalIp()
	if err != nil {
		return err
	}
	nodeOpts := &node.Node{
		Name:        configs.Default.NodeName,
		Version:     configs.Version,
		CreateTime:  time.Now(),
		Role:        node.Role{Master: true, Admin: true, Worker: true},
		LastRunTime: time.Now(),
		IP:          nip.String(),
	}
	if _, err = dao.RefreshNode(c, nodeOpts); err != nil {
		return
	}
	if _, err = dao.CreateStatusIfNotExist(c, nodeOpts.Name, node.WithReady(true)); err != nil {
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
			logrus.WithField("node name", configs.GetNodeName()).Info("Running node")
			err = server.Run(ctx)
		}()
		time.Sleep(time.Second * 5)
	}
}
