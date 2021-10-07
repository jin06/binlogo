package api

import (
	"github.com/gin-gonic/gin"
	routers2 "github.com/jin06/binlogo/app/server/api/routers"
)

type APIServer struct {

}

func NewAPIServer() *APIServer {
	return &APIServer{}
}

func (server *APIServer) Run() {
	g := gin.Default()
	// default 8080

	g.GET("/pipelines", routers2.List)
	g.GET("/pipelines/:id", routers2.One)
	g.POST("/pipelines", routers2.Create)
	g.PUT("/pipelines/:id", routers2.Update)
	g.DELETE("/pipelines/:id", routers2.Delete)
	g.Run()

}
