package instance

import (
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/app/server/console/handler"
	"github.com/jin06/binlogo/app/server/console/module/instance"
)

func Get(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(200, handler.Fail("Name is null"))
		return
	}
	item, err := instance.GetInstanceByName(name)
	if err != nil {
		c.JSON(200, handler.Fail(err))
		return
	}
	c.JSON(200, handler.Success(item))
}
