package pipeline

import (
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/app/server/console/handler"
	"github.com/jin06/binlogo/app/server/console/module/pipeline"
)

func Get(c *gin.Context) {
	q := struct {
		Name string `json:"name"`
	}{}
	if err := c.Bind(&q); err != nil {
		c.JSON(200, handler.Fail(err))
	}
	item , err := pipeline.GetItemByName(q.Name)
	if err != nil {
		c.JSON(200, handler.Fail(err))
	}

	c.JSON(200, handler.Success(item))
}
