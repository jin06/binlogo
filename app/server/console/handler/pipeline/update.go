package pipeline

import (
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/v2/app/server/console/handler"
	pipeline2 "github.com/jin06/binlogo/v2/app/server/console/module/pipeline"
	"github.com/jin06/binlogo/v2/pkg/pipeline/tool"
	"github.com/jin06/binlogo/v2/pkg/store/dao/dao_pipe"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
)

// Update handler, update pipline
func Update(c *gin.Context) {
	q := &pipeline.Pipeline{}
	if err := c.BindJSON(q); err != nil {
		c.JSON(200, handler.Fail(err.Error()))
		return
	}
	for _, v := range q.Filters {
		if !tool.FilterVerifyStr(v.Rule) {
			c.JSON(200, handler.Fail("Filter rule error, only support the format like database.table or database "))
			return
		}
	}

	pipe, err := dao_pipe.GetPipeline(q.Name)
	if err != nil {
		c.JSON(200, handler.Fail(err.Error()))
		return
	}
	if pipe.Status == pipeline.STATUS_RUN {
		c.JSON(200, handler.Fail("Only stopped pipeline can be updated"))
		return
	}
	pipelineDefault(q)
	ok, err := dao_pipe.UpdatePipeline(q.Name, pipeline.WithPipeSafe(q))
	if err != nil || !ok {
		c.JSON(200, "update failed")
		return
	}
	item, _ := pipeline2.GetItemByName(q.Name)

	c.JSON(200, handler.Success(item))
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
	//fmt.Println(q)
	if q.Status != pipeline.STATUS_RUN && q.Status != pipeline.STATUS_STOP {
		c.JSON(200, handler.Fail("Wrong param status: "+q.Status))
		return
	}

	ok, err := dao_pipe.UpdatePipeline(q.PipeName, pipeline.WithPipeStatus(q.Status))
	if err != nil || !ok {
		c.JSON(200, handler.Fail("Update status failed "))
		return
	}
	c.JSON(200, handler.Success("ok"))
}

func UpdateMode(c *gin.Context) {
	q := struct {
		PipeName string        `json:"name"`
		Mode     pipeline.Mode `json:"mode"`
	}{}
	if err := c.BindJSON(&q); err != nil {
		c.JSON(200, handler.Fail(err))
		return
	}
	if q.Mode != pipeline.MODE_POSITION && q.Mode != pipeline.MODE_GTID {
		c.JSON(200, handler.Fail("Wrong param mode: "+q.Mode))
		return
	}
	pipe, err := dao_pipe.GetPipeline(q.PipeName)
	if err != nil {
		c.JSON(200, handler.Fail(err.Error()))
		return
	}
	if pipe.Status == pipeline.STATUS_RUN {
		c.JSON(200, handler.Fail("Only stopped pipeline can be updated"))
		return
	}
	ok, err := dao_pipe.UpdatePipeline(q.PipeName, pipeline.WithPipeMode(q.Mode))
	if err != nil || !ok {
		c.JSON(200, handler.Fail("Update mode failed"))
		return
	}
	c.JSON(200, handler.Success("ok"))

}
