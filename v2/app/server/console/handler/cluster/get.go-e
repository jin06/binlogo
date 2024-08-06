package cluster

import (
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/app/server/console/handler"
	"github.com/spf13/viper"
)

func Get(c *gin.Context) {
	data := &struct {
		Name string `json:"name"`
	}{}

	data.Name = viper.GetString("cluster.name")

	c.JSON(200, handler.Success(data))
	return
}
