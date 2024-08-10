package pipeline

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/v2/app/server/console/handler"
	"github.com/jin06/binlogo/v2/pkg/pipeline/tool"
	"github.com/jin06/binlogo/v2/pkg/store/dao/dao_pipe"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
	"github.com/jin06/binlogo/v2/pkg/util/random"
	"github.com/sirupsen/logrus"
)

func Create(c *gin.Context) {
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
	q.CreateTime = time.Now()

	logrus.Debugf("%v \n", *q)
	if q.Mysql.ServerId == 0 {
		//q.Mysql.ServerId = uint32(rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(100000000))
		q.Mysql.ServerId = random.Uint32()
	}
	q.Status = pipeline.STATUS_STOP
	pipelineDefault(q)
	if _, err := dao_pipe.CreatePipeline(q); err != nil {
		c.JSON(200, handler.Fail(err))
		return
	}
	c.JSON(200, handler.Success("ok"))
}

func pipelineDefault(p *pipeline.Pipeline) {
	switch p.Output.Sender.Type {
	case pipeline.SNEDER_TYPE_RABBITMQ:
		{
			if p.Output.Sender.RabbitMQ.ExchangeName == "" {
				p.Output.Sender.RabbitMQ.ExchangeName = p.Name
			}
		}
	case pipeline.SENDER_TYPE_REDIS:
		{
			if p.Output.Sender.Redis.List == "" {
				p.Output.Sender.Redis.List = p.Name
			}
		}
	case pipeline.SENDER_TYPE_KAFKA:
		{
			if p.Output.Sender.Kafka.Topic == "" {
				p.Output.Sender.Kafka.Topic = p.Name
			}
		}
	}
}
