package position

import (
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/app/server/console/handler"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
)

// Get handler for get pipeline by pipeline name
func Get(c *gin.Context) {
	name := c.Query("pipe_name")
	if name == "" {
		c.JSON(200, handler.Fail("pipe_name is empty"))
		return
	}
	pos, err := dao_pipe.GetPosition(name)
	if err != nil {
		c.JSON(200, handler.Fail(err))
		return
	}
	c.JSON(200, handler.Success(pos))
	return
}
