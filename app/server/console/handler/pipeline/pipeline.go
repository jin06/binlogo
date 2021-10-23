package pipeline

import (
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/app/server/console/handler"
)

func Get(c *gin.Context) {
	c.JSON(200, handler.Success(nil))
}

func Create(c *gin.Context) {
	c.JSON(200, handler.Success(nil))
}

func Update(c *gin.Context) {
	c.JSON(200, handler.Success(nil))
}

func Delete(c *gin.Context) {
	c.JSON(200, handler.Success(nil))
}
