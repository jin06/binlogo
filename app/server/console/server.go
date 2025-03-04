package console

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/v2/configs"
	"github.com/sirupsen/logrus"
)

func Run(ctx context.Context) (err error) {
	g := gin.Default()
	router(g)
	listen := fmt.Sprintf("%s:%d", configs.Default.Console.Listen, configs.Default.Console.Port)
	logrus.Info("Console started, listen ==>  ", listen)
	return g.Run(listen)
}
