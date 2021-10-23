package console

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/app/server/console/handler/pipeline"
	"github.com/jin06/binlogo/app/server/console/middleware"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/spf13/viper"
)

func Run(ctx context.Context) (err error) {
	g := gin.Default()
	g.Use(middleware.Cors)

	g.GET("/api/pipeline/list", pipeline.List)
	g.GET("/api/pipeline/get", pipeline.Get)
	g.POST("/api/pipeline/create", pipeline.Create)
	g.POST("/api/pipeline/update", pipeline.Update)
	g.POST("/api/pipeline/delete", pipeline.Delete)
	listen := viper.GetString("console.listen") + ":" + viper.GetString("console.port")
	blog.Info("Console api --> ", listen)
	err = g.Run(listen)
	return
}
