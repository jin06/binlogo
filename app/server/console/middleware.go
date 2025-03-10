package console

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/v2/app/server/console/basic"
	"github.com/jin06/binlogo/v2/app/server/console/service"
)

func corsMiddle() gin.HandlerFunc {
	cfg := cors.DefaultConfig()
	cfg.AddAllowHeaders("x-token")
	cfg.AllowAllOrigins = true
	return cors.New(cfg)
}

func auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("x-token")
		if !service.DefaultStore().Get(token) {
			c.AbortWithStatusJSON(200, basic.FailCode(basic.CodeTokenExpired))
		}
	}
}
