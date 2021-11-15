package console

import (
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/app/server/console/handler/instance"
	"github.com/jin06/binlogo/app/server/console/handler/node"
	"github.com/jin06/binlogo/app/server/console/handler/pipeline"
	"github.com/jin06/binlogo/app/server/console/middleware"
)

func router(g *gin.Engine) {
	g.Use(middleware.Cors)

	g.GET("/api/pipeline/list", pipeline.List)
	g.GET("/api/pipeline/get", pipeline.Get)
	g.POST("/api/pipeline/create", pipeline.Create)
	g.POST("/api/pipeline/update", pipeline.Update)
	g.POST("/api/pipeline/update/status", pipeline.UpdateStatus)
	g.POST("/api/pipeline/delete", pipeline.Delete)

	g.GET("/api/node/list", node.List)

	g.GET("/api/instance/get", instance.Get)
}
