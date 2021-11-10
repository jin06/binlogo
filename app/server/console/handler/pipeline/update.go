package pipeline

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/app/server/console/handler"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
)

func Update(c *gin.Context) {
	c.JSON(200, handler.Success("ok"))
}

func UpdateStatus(c *gin.Context) {
	q := struct {
		PipeName string          `json:"name"`
		Status   pipeline.Status `json:"status"`
	}{}
	if err := c.BindJSON(&q); err != nil {
		c.JSON(200, handler.Fail(err.Error()))
		return
	}
	fmt.Println(q)
	if q.Status != pipeline.STATUS_RUN && q.Status != pipeline.STATUS_STOP {
		c.JSON(200, handler.Fail("wrong param status: "+q.Status))
		return
	}
	ok, err := dao_pipe.UpdatePipeline(q.PipeName, pipeline.WithPipeStatus(q.Status))
	if err != nil || !ok {
		c.JSON(200, handler.Fail("udpate status failed "))
		return
	}
	c.JSON(200, handler.Success("ok"))
}
