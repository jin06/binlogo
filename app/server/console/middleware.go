package console

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func corsMiddle() gin.HandlerFunc {
	cfg := cors.DefaultConfig()
	cfg.AddAllowHeaders("x-token")
	cfg.AllowAllOrigins = true
	return cors.New(cfg)
}
