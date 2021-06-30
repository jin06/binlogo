package api

import "github.com/gin-gonic/gin"

type Server struct {

}

func NewApiServer() *Server {
	return &Server{}
}

func (server *Server) Run() {
	router := gin.Default()
	// default 8080
	router.Run()
}