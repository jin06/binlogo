package console

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Run(ctx context.Context) (err error) {
	g := gin.Default()

	router(g)

	listen := viper.GetString("console.listen") + ":" + viper.GetString("console.port")
	logrus.Info("Console api --> ", listen)
	err = g.Run(listen)
	return
}
