package pipeline

import (
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/app/server/console/handler"
	"github.com/jin06/binlogo/app/server/console/module/pipeline"
)

func Get(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(200,handler.Fail("Name is null"))
	}

	item , err := pipeline.GetItemByName(name)
	if err != nil {
		c.JSON(200, handler.Fail(err))
	}

	c.JSON(200, handler.Success(item))
}
