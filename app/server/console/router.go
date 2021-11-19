package console

import (
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/app/server/console/handler/cluster"
	"github.com/jin06/binlogo/app/server/console/handler/event"
	"github.com/jin06/binlogo/app/server/console/handler/instance"
	"github.com/jin06/binlogo/app/server/console/handler/node"
	"github.com/jin06/binlogo/app/server/console/handler/pipeline"
	"github.com/jin06/binlogo/app/server/console/handler/position"
	"github.com/jin06/binlogo/app/server/console/middleware"
	"net/http"
)

func router(g *gin.Engine) {
	g.Use(middleware.Cors)

	g.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/console")
	})
	g.Static("/console", "./assets/dist")

	g.GET("/api/pipeline/list", pipeline.List)
	g.GET("/api/pipeline/get", pipeline.Get)
	g.POST("/api/pipeline/create", pipeline.Create)
	g.POST("/api/pipeline/update", pipeline.Update)
	g.POST("/api/pipeline/update/status", pipeline.UpdateStatus)
	g.POST("/api/pipeline/update/mode", pipeline.UpdateMode)
	g.POST("/api/pipeline/delete", pipeline.Delete)
	g.GET("/api/pipeline/is_filter", pipeline.IsFilter)
	g.POST("/api/pipeline/add_filter", pipeline.AddFilter)
	g.POST("/api/pipeline/update_filter", pipeline.UpdateFilter)

	g.GET("/api/node/list", node.List)

	g.GET("/api/instance/get", instance.Get)
	g.GET("/api/instance/list", instance.List)

	g.GET("/api/position/get", position.Get)
	g.POST("/api/position/update", position.Update)

	g.GET("/api/cluster/get", cluster.Get)
	g.GET("/api/cluster/list/register", cluster.RegisterList)
	g.GET("/api/cluster/list/election", cluster.ElectionList)

	g.GET("/api/event/list", event.List)
}
