package console

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/v2/app/server/console/handler"
	mid "github.com/jin06/binlogo/v2/app/server/console/middleware"
	"github.com/jin06/binlogo/v2/static/dist"
)

func router(g *gin.Engine) {
	g.Use(corsMiddle())
	// g.GET("/metrics", gin.WrapH(promhttp.Handler()))
	g.Use(mid.Cors)

	staticServer := http.FS(dist.Content)
	g.StaticFS("/console", staticServer)
	// g.StaticFS("/assets")
	// g.Static("/assets", "./assets/dist/assets")

	g.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/console")
		// c.FileFromFS("console/index.html", staticServer)
	})

	g.POST("/api/user/login", handler.Login)
	g.POST("/api/user/logout", handler.Logout)
	g.GET("/api/user/auth_type", handler.AuthType)

	g.Use(auth())

	g.GET("/api/user/info", handler.UserInfo)

	g.GET("/api/pipeline/list", handler.PipeList)
	g.GET("/api/pipeline/get", handler.PipeGet)
	g.POST("/api/pipeline/create", handler.PipeCreate)
	g.POST("/api/pipeline/update", handler.PipeUpdate)
	g.POST("/api/pipeline/update/status", handler.UpdateStatus)
	g.POST("/api/pipeline/update/mode", handler.UpdateMode)
	g.POST("/api/pipeline/delete", handler.PipeDelete)
	g.GET("/api/pipeline/is_filter", handler.PipeIsFilter)
	g.POST("/api/pipeline/add_filter", handler.PipeAddFilter)
	g.POST("/api/pipeline/update_filter", handler.PipeUpdateFilter)

	g.GET("/api/node/list", handler.NodeList)

	g.GET("/api/instance/get", handler.InstanceGet)
	g.GET("/api/instance/list", handler.InstanceList)

	g.GET("/api/position/get", handler.PositionGet)
	g.POST("/api/position/update", handler.PositionUpdate)

	g.GET("/api/cluster/get", handler.ClusterGet)
	g.GET("/api/cluster/list/register", handler.RegisterList)
	g.GET("/api/cluster/list/election", handler.ElectionList)

	g.GET("/api/event/list", handler.EventList)
	g.GET("/api/event/scroll_list", handler.EventScrollList)

}
