package app

import (
	"context"
	"fmt"
	"time"

	server_node "github.com/jin06/binlogo/v2/app/server/node"
	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/pkg/blog"
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
	dao.CreateStatusIfNotExist(c, nodeOpts.Name, node.StatusConditions{
		node.ConReady: true,
	})

	logrus.Info("run node")
	event.Event(model_event.NewInfoNode("Run node success"))
	for i := 1; ; i++ {
		func(i int) {
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
			// log := logrus.WithField("NodeSeq", i)
			log := blog.NewLog().WithField("NodeSeq", i)
			server := server_node.New(&server_node.Options{Node: nodeOpts, Log: log})
			log.Info("Running node!")
			// logrus.WithField("name", configs.GetNodeName()).Info("Running node")
			err = server.Run(ctx)
		}(i)
		time.Sleep(time.Second * 5)
	}
}
