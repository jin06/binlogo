package position

import (
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/v2/app/server/console/handler"
	pipeline2 "github.com/jin06/binlogo/v2/app/server/console/module/pipeline"
	"github.com/jin06/binlogo/v2/pkg/store/dao"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
)

func Update(c *gin.Context) {
	q := struct {
		Mode     pipeline.Mode      `json:"mode"`
		Position *pipeline.Position `json:"position"`
	}{}
	if err := c.BindJSON(&q); err != nil {
		c.JSON(200, handler.Fail(err))
		return
	}
	if q.Mode != pipeline.MODE_POSITION && q.Mode != pipeline.MODE_GTID {
		c.JSON(200, handler.Fail("mode error: "+q.Mode))
		return
	}
	pipeStatus, err := pipeline2.PipeStatus(q.Position.PipelineName)
	if err != nil {
		c.JSON(200, handler.Fail(err))
		return
	}
	if pipeStatus == pipeline.STATUS_RUN {
		c.JSON(200, handler.Fail("Only stopped pipeline can be updated"))
		return
	}
	switch q.Mode {
	case pipeline.MODE_GTID:
		{
			var ok bool
			ok, err = dao.UpdateRecordSafe(
				q.Position.PipelineName,
				pipeline.WithPre(&pipeline.Position{
					PipelineName: q.Position.PipelineName,
					GTIDSet:      q.Position.GTIDSet,
				}),
				pipeline.WithNow(nil),
			)
			if err != nil || !ok {
				c.JSON(200, handler.Fail("update failed "+err.Error()))
				return
			}
		}
	case pipeline.MODE_POSITION:
		{
			var ok bool
			ok, err = dao.UpdateRecordSafe(
				q.Position.PipelineName,
				pipeline.WithPre(&pipeline.Position{
					BinlogFile:     q.Position.BinlogFile,
					BinlogPosition: q.Position.BinlogPosition,
					PipelineName:   q.Position.PipelineName,
				}),
				pipeline.WithNow(nil),
			)
			if err != nil || !ok {
				c.JSON(200, handler.Fail("update failed "+err.Error()))
				return
			}
		}
	default:
		c.JSON(200, handler.Fail("mode error: "+q.Mode))
		return
	}
	c.JSON(200, handler.Success("ok"))
}
