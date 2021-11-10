package pipeline

import (
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/app/server/console/handler"
	"github.com/jin06/binlogo/app/server/console/module/pipeline"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
)

func Get(c *gin.Context) {
	q := struct {
		Name string `json:"name"`
	}{}
	if err := c.Bind(&q); err != nil {
		c.JSON(200, handler.Fail(err))
	}
	pipe, err := dao_pipe.GetPipeline(q.Name)
	if err != nil {
		c.JSON(200, handler.Fail(err))
		return
	}
	item := &pipeline.Item{
		Pipeline: pipe,
	}
	err = pipeline.CompleteInfo(item)
	if err != nil {
		c.JSON(200, handler.Fail(err))
		return
	}
	c.JSON(200, handler.Success(nil))
}
