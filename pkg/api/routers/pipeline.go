package routers

import "github.com/gin-gonic/gin"

func List(c *gin.Context) {
	c.JSON(200, NewSuccess())
}

func One(c *gin.Context) {
	c.JSON(200, NewSuccess())
}

func Create(c *gin.Context) {
	c.JSON(200, NewSuccess())
}

func Update(c *gin.Context)  {
	c.JSON(200, NewSuccess())
}

func Delete(c *gin.Context)  {
	c.JSON(200, NewSuccess())
}