package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors(c *gin.Context) {
	method := c.Request.Method
	origin := c.Request.Header.Get("Origin") //请求头部
	if origin != "" {
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "172800")
	}

	if method == "OPTIONS" {
		if c.GetHeader("Access-Control-Request-Method") != "" {
			c.Header("Access-Control-Allow-Methods", c.GetHeader("Access-Control-Request-Method"))
		}
		if c.GetHeader("Access-Control-Request-Headers") != "" {
			c.Header("Access-Control-Allow-Headers", c.GetHeader("Access-Control-Request-Headers"))
		}
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
}
