package pipeline

import (
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/v2/app/server/console/handler"
	"github.com/jin06/binlogo/v2/pkg/store/dao"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
)

func Delete(c *gin.Context) {
	q := struct {
		Name string `json:"name"`
	}{}
	if err := c.BindJSON(&q); err != nil {
		c.JSON(200, handler.Fail(err.Error()))
		return
	}
	pipe, err := dao.GetPipeline(q.Name)
	if err != nil {
		c.JSON(200, handler.Fail(err.Error()))
		return
	}
	if pipe.Status == pipeline.STATUS_RUN {
		c.JSON(200, handler.Fail("Only stopped pipeline can be deleted"))
		return
	}
	ok, err := dao.UpdatePipeline(q.Name, pipeline.WithPipeDelete(true))
	if err != nil || !ok {
		c.JSON(200, handler.Fail("Delete pipeline failed"))
		return
	}

	c.JSON(200, handler.Success("ok"))
}
