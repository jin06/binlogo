package pipeline

import (
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/app/server/console/handler"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"time"
)

func Create(c *gin.Context) {
	q := &pipeline.Pipeline{}

	if err := c.BindJSON(q); err != nil {
		c.JSON(200, handler.Fail(err.Error()))
		return
	}
	q.CreateTime = time.Now()

	blog.Debugf("%v \n", *q)
	if _, err := dao_pipe.CreatePipeline(q) ; err != nil {
		c.JSON(200, handler.Fail(err))
		return
	}
	c.JSON(200, handler.Success("ok"))
}
