package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/v2/app/server/console/basic"
	pipeModule "github.com/jin06/binlogo/v2/app/server/console/module/pipeline"
	"github.com/jin06/binlogo/v2/pkg/store/dao"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
)

// Get handler for get pipeline by pipeline name
func PositionGet(c *gin.Context) {
	name := c.Query("pipe_name")
	if name == "" {
		c.JSON(200, basic.Fail("pipe_name is empty"))
		return
	}
	//pos, err := dao.GetPosition(name)
	record, err := dao.GetRecord(c, name)
	if err != nil {
		c.JSON(200, basic.Fail(err))
		return
	}
	res := &pipeline.Position{}
	if record != nil {
		if record.Pre != nil {
			res = record.Pre
		}
	}
	c.JSON(200, basic.Success(res))
}

func PositionUpdate(c *gin.Context) {
	q := struct {
		Mode     pipeline.Mode      `json:"mode"`
		Position *pipeline.Position `json:"position"`
	}{}
	if err := c.BindJSON(&q); err != nil {
		c.JSON(200, basic.Fail(err))
		return
	}
	if q.Mode != pipeline.MODE_POSITION && q.Mode != pipeline.MODE_GTID {
		c.JSON(200, basic.Fail("mode error: "+q.Mode))
		return
	}
	pipeStatus, err := pipeModule.PipeStatus(c, q.Position.PipelineName)
	if err != nil {
		c.JSON(200, basic.Fail(err))
		return
	}
	if pipeStatus == pipeline.STATUS_RUN {
		c.JSON(200, basic.Fail("Only stopped pipeline can be updated"))
		return
	}
	switch q.Mode {
	case pipeline.MODE_GTID:
		{
			err = dao.UpdateRecordSafe(
				c,
				q.Position.PipelineName,
				pipeline.WithPre(&pipeline.Position{
					PipelineName: q.Position.PipelineName,
					GTIDSet:      q.Position.GTIDSet,
				}),
				pipeline.WithNow(nil),
			)
			if err != nil {
				c.JSON(200, basic.Fail("update failed "+err.Error()))
				return
			}
		}
	case pipeline.MODE_POSITION:
		{
			err = dao.UpdateRecordSafe(
				c,
				q.Position.PipelineName,
				pipeline.WithPre(&pipeline.Position{
					BinlogFile:     q.Position.BinlogFile,
					BinlogPosition: q.Position.BinlogPosition,
					PipelineName:   q.Position.PipelineName,
				}),
				pipeline.WithNow(nil),
			)
			if err != nil {
				c.JSON(200, basic.Fail("update failed "+err.Error()))
				return
			}
		}
	default:
		c.JSON(200, basic.Fail("mode error: "+q.Mode))
		return
	}
	c.JSON(200, basic.Success("ok"))
}
