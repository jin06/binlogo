package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/server/api/routers"
)

type APIServer struct {

}

func NewAPIServer() *APIServer {
	return &APIServer{}
}

func (server *APIServer) Run() {
	g := gin.Default()
	// default 8080

	g.GET("/pipelines", routers.List)
	g.GET("/pipelines/:id",routers.One)
	g.POST("/pipelines", routers.Create)
	g.PUT("/pipelines/:id", routers.Update)
	g.DELETE("/pipelines/:id", routers.Delete)
	g.Run()

}
