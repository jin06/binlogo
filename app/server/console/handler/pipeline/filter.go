package pipeline

import (
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/app/server/console/handler"
	"github.com/jin06/binlogo/pkg/pipeline/pipe_tool"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
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
	if err != nil {
		c.JSON(200, handler.Fail(err))
		return
	}
	c.JSON(200, handler.Success(isFilter))
	return
}

func AddFilter(c *gin.Context) {
	q := &struct {
		PipeName string           `json:"pipe_name"`
		Filter   *pipeline.Filter `json:"filter"`
	}{}
	if err := c.BindJSON(q);err != nil {
		c.JSON(200, handler.Fail(err))
		return
	}
	ok, err := dao_pipe.UpdatePipeline(q.PipeName, pipeline.WithAddFilter(q.Filter))
	if err != nil  || !ok {
		c.JSON(200, handler.Fail("Add filter failed."))
		return
	}
	c.JSON(200, handler.Success("ok"))
}

func UpdateFilter(c *gin.Context) {
	q := &struct {
		PipeName string           `json:"pipe_name"`
		Index    int              `json:"index"`
		Filter   *pipeline.Filter `json:"filter"`
	}{}
	if err := c.BindJSON(q);err != nil {
		c.JSON(200, handler.Fail(err))
		return
	}
	ok, err := dao_pipe.UpdatePipeline(q.PipeName, pipeline.WithUpdateFilter(q.Index, q.Filter))
	if err != nil  || !ok {
		c.JSON(200, handler.Fail("Update filter failed."))
		return
	}
	c.JSON(200, handler.Success("ok"))
}
