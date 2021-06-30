package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/pkg/api/routers"
)

type Server struct {

}

func NewApiServer() *Server {
	return &Server{}
}

func (server *Server) Run() {
	g := gin.Default()
	// default 8080

	g.GET("/pipelines", routers.List)
	g.GET("/pipelines/:id",routers.One)
	g.POST("/pipelines", routers.Create)
	g.PUT("/pipelines/:id", routers.Update)
	g.DELETE("/pipelines/:id", routers.Delete)
	g.Run()

}