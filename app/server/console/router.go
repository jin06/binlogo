package console

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/v2/app/server/console/handler/cluster"
	"github.com/jin06/binlogo/v2/app/server/console/handler/event"
	"github.com/jin06/binlogo/v2/app/server/console/handler/instance"
	"github.com/jin06/binlogo/v2/app/server/console/handler/node"
	"github.com/jin06/binlogo/v2/app/server/console/handler/pipeline"
	"github.com/jin06/binlogo/v2/app/server/console/handler/position"
	"github.com/jin06/binlogo/v2/app/server/console/handler/user"
	mid "github.com/jin06/binlogo/v2/app/server/console/middleware"
	"github.com/jin06/binlogo/v2/static"
)

func router(g *gin.Engine) {
	g.Use(corsMiddle())
	// g.GET("/metrics", gin.WrapH(promhttp.Handler()))
	g.Use(mid.Cors)

	staticServer := http.FS(static.Content)
	g.StaticFS("/console", staticServer)
	// g.StaticFS("/assets")
	// g.Static("/assets", "./assets/dist/assets")

	g.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/console")
		// c.FileFromFS("console/index.html", staticServer)
	})

	g.POST("/api/user/login", user.Login)
	g.POST("/api/user/logout", user.Logout)
	g.GET("/api/user/auth_type", user.AuthType)

	g.Use(auth())

	g.GET("/api/user/info", user.Info)

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
	g.GET("/api/event/scroll_list", event.ScrollList)

}
