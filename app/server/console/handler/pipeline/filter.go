package pipeline

import (
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/app/server/console/handler"
	"github.com/jin06/binlogo/pkg/pipeline/pipe_tool"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
)

func IsFilter(c *gin.Context) {
	rule := c.Query("rule")
	name := c.Query("name")
	if rule == "" {
		c.JSON(200, handler.Fail("params is null"))
		return
	}
	if name == "" {
		c.JSON(200, handler.Fail("fatal error, refresh page and try again"))
		return
	}
	pipe, err := dao_pipe.GetPipeline(name)
	if err != nil {
		c.JSON(200, handler.Fail(err))
		return
	}
	validator := pipe_tool.NewFilter(pipe.Filters)
	isFilter, err := validator.IsFilterWithName(rule)
	if err !=nil {
		c.JSON(200, handler.Fail(err))
		return
	}
	c.JSON(200, handler.Success(isFilter))
	return
}
