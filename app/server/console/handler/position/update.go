package position

import (
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/app/server/console/handler"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
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
	if q.Mode != pipeline.MODE_POSTION && q.Mode != pipeline.MODE_GTID {
		c.JSON(200, handler.Fail("mode error: "+q.Mode))
		return
	}
	switch q.Mode {
	case pipeline.MODE_GTID:
		{
			ok, err := dao_pipe.UpdatePositionSafe(q.Position.PipelineName, pipeline.WithGTIDSet(q.Position.GTIDSet))
			if err != nil || !ok {
				c.JSON(200, handler.Fail("update failed " + err.Error()))
				return
			}
		}
	case pipeline.MODE_POSTION:
		{
			ok, err := dao_pipe.UpdatePositionSafe(q.Position.PipelineName, pipeline.WithBinlogFile(q.Position.BinlogFile), pipeline.WithPos(q.Position.BinlogPosition))
			if err != nil || !ok {
				c.JSON(200, handler.Fail("update failed " + err.Error()))
				return
			}
		}
	default:
		c.JSON(200, handler.Fail("mode error: "+q.Mode))
		return
	}
	c.JSON(200, handler.Success("ok"))
}
