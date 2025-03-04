package handler

import (
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/v2/app/server/console/basic"
	pipeModule "github.com/jin06/binlogo/v2/app/server/console/module/pipeline"
	"github.com/jin06/binlogo/v2/pkg/pipeline/tool"
	"github.com/jin06/binlogo/v2/pkg/store/dao"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
	"github.com/jin06/binlogo/v2/pkg/util/random"
	"github.com/sirupsen/logrus"
)

func PipeList(c *gin.Context) {
	qSort := c.Query("sort")
	name := c.Query("name")
	status := c.Query("status")

	all, err := dao.AllPipelines(c)
	if err != nil {
		c.JSON(200, basic.Fail(err))
		return
	}
	var items []*pipeModule.Item
	for _, v := range all {
		if v.IsDelete {
			continue
		}
		items = append(items, &pipeModule.Item{Pipeline: v})
	}

	if err = pipeModule.CompleteInfoList(c, items); err != nil {
		c.JSON(200, basic.Fail(err))
		return
	}

	resItems := []*pipeModule.Item{}
	for _, v := range items {
		if status != "" {
			if string(v.Pipeline.Status) != status {
				continue
			}
		}
		if name != "" {
			if !strings.Contains(v.Pipeline.Name, name) && !strings.Contains(v.Pipeline.AliasName, name) {
				continue
			}
		}
		resItems = append(resItems, v)
	}

	sort.Slice(resItems, func(i, j int) bool {
		if qSort == "+id" {
			return resItems[i].Pipeline.CreateTime.Before(resItems[j].Pipeline.CreateTime)
		} else {
			return resItems[j].Pipeline.CreateTime.Before(resItems[i].Pipeline.CreateTime)
		}
	})
	start, end := basic.StartEnd(c)

	if end > len(resItems) {
		end = len(resItems)
	}

	c.JSON(200, basic.Success(map[string]interface{}{
		"items": resItems[start:end],
		"total": len(resItems),
	}))
}

func PipeGet(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(200, basic.Fail("Name is null"))
		return
	}

	item, err := pipeModule.GetItemByName(c, name)
	if err != nil {
		c.JSON(200, basic.Fail(err))
		return
	}

	c.JSON(200, basic.Success(item))
}

func PipeCreate(c *gin.Context) {
	q := &pipeline.Pipeline{}

	if err := c.BindJSON(q); err != nil {
		c.JSON(200, basic.Fail(err.Error()))
		return
	}
	for _, v := range q.Filters {
		if !tool.FilterVerifyStr(v.Rule) {
			c.JSON(200, basic.Fail("Filter rule error, only support the format like database.table or database "))
			return
		}
	}
	q.CreateTime = time.Now()

	logrus.Debugf("%v \n", *q)
	if q.Mysql.ServerId == 0 {
		q.Mysql.ServerId = random.Uint32()
	}
	q.Status = pipeline.STATUS_STOP
	pipelineDefault(q)
	if _, err := dao.CreatePipeline(c, q); err != nil {
		c.JSON(200, basic.Fail(err))
		return
	}
	c.JSON(200, basic.Success("ok"))
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

func PipeDelete(c *gin.Context) {
	q := struct {
		Name string `json:"name"`
	}{}
	if err := c.BindJSON(&q); err != nil {
		c.JSON(200, basic.Fail(err.Error()))
		return
	}
	pipe, err := dao.GetPipeline(c, q.Name)
	if err != nil {
		c.JSON(200, basic.Fail(err.Error()))
		return
	}
	if pipe.Status == pipeline.STATUS_RUN {
		c.JSON(200, basic.Fail("Only stopped pipeline can be deleted"))
		return
	}
	err = dao.UpdatePipeline(c, q.Name, pipeline.WithPipeDelete(true))
	if err != nil {
		c.JSON(200, basic.Fail("Delete pipeline failed"))
		return
	}

	c.JSON(200, basic.Success("ok"))
}

func PipeIsFilter(c *gin.Context) {
	rule := c.Query("rule")
	name := c.Query("name")
	if rule == "" {
		c.JSON(200, basic.Fail("params is null"))
		return
	}
	if name == "" {
		c.JSON(200, basic.Fail("fatal error, refresh page and try again"))
		return
	}
	pipe, err := dao.GetPipeline(c, name)
	if err != nil {
		c.JSON(200, basic.Fail(err))
		return
	}
	validator := tool.NewFilter(pipe.Filters)
	isFilter, err := validator.IsFilterWithName(rule)
	if err != nil {
		c.JSON(200, basic.Fail(err))
		return
	}
	c.JSON(200, basic.Success(isFilter))
	return
}

func PipeAddFilter(c *gin.Context) {
	q := &struct {
		PipeName string           `json:"pipe_name"`
		Filter   *pipeline.Filter `json:"filter"`
	}{}
	if err := c.BindJSON(q); err != nil {
		c.JSON(200, basic.Fail(err))
		return
	}
	err := dao.UpdatePipeline(c, q.PipeName, pipeline.WithAddFilter(q.Filter))
	if err != nil {
		c.JSON(200, basic.Fail("Add filter failed."))
		return
	}
	c.JSON(200, basic.Success("ok"))
}

func PipeUpdateFilter(c *gin.Context) {
	q := &struct {
		PipeName string           `json:"pipe_name"`
		Index    int              `json:"index"`
		Filter   *pipeline.Filter `json:"filter"`
	}{}
	if err := c.BindJSON(q); err != nil {
		c.JSON(200, basic.Fail(err))
		return
	}
	err := dao.UpdatePipeline(c, q.PipeName, pipeline.WithUpdateFilter(q.Index, q.Filter))
	if err != nil {
		c.JSON(200, basic.Fail("Update filter failed."))
		return
	}
	c.JSON(200, basic.Success("ok"))
}

// Update handler, update pipline
func PipeUpdate(c *gin.Context) {
	q := &pipeline.Pipeline{}
	if err := c.BindJSON(q); err != nil {
		c.JSON(200, basic.Fail(err.Error()))
		return
	}
	for _, v := range q.Filters {
		if !tool.FilterVerifyStr(v.Rule) {
			c.JSON(200, basic.Fail("Filter rule error, only support the format like database.table or database "))
			return
		}
	}

	pipe, err := dao.GetPipeline(c, q.Name)
	if err != nil {
		c.JSON(200, basic.Fail(err.Error()))
		return
	}
	if pipe.Status == pipeline.STATUS_RUN {
		c.JSON(200, basic.Fail("Only stopped pipeline can be updated"))
		return
	}
	pipelineDefault(q)
	if err = dao.UpdatePipeline(c, q.Name, pipeline.WithPipeSafe(q)); err != nil {
		c.JSON(200, "update failed")
		return
	}
	item, _ := pipeModule.GetItemByName(c, q.Name)

	c.JSON(200, basic.Success(item))
}

type updateStatusReq struct {
	PipeName string          `json:"name"`
	Status   pipeline.Status `json:"status"`
}

func UpdateStatus(c *gin.Context) {
	q := updateStatusReq{}
	if err := c.BindJSON(&q); err != nil {
		c.JSON(200, basic.Fail(err.Error()))
		return
	}
	//fmt.Println(q)
	if q.Status != pipeline.STATUS_RUN && q.Status != pipeline.STATUS_STOP {
		c.JSON(200, basic.Fail("Wrong param status: "+q.Status))
		return
	}

	err := dao.UpdatePipeline(c, q.PipeName, pipeline.WithPipeStatus(q.Status))
	if err != nil {
		c.JSON(200, basic.Fail("Update status failed "))
		return
	}
	c.JSON(200, basic.Success("ok"))
}

func UpdateMode(c *gin.Context) {
	q := struct {
		PipeName string        `json:"name"`
		Mode     pipeline.Mode `json:"mode"`
	}{}
	if err := c.BindJSON(&q); err != nil {
		c.JSON(200, basic.Fail(err))
		return
	}
	if q.Mode != pipeline.MODE_POSITION && q.Mode != pipeline.MODE_GTID {
		c.JSON(200, basic.Fail("Wrong param mode: "+q.Mode))
		return
	}
	pipe, err := dao.GetPipeline(c, q.PipeName)
	if err != nil {
		c.JSON(200, basic.Fail(err.Error()))
		return
	}
	if pipe.Status == pipeline.STATUS_RUN {
		c.JSON(200, basic.Fail("Only stopped pipeline can be updated"))
		return
	}
	err = dao.UpdatePipeline(c, q.PipeName, pipeline.WithPipeMode(q.Mode))
	if err != nil {
		c.JSON(200, basic.Fail("Update mode failed"))
		return
	}
	c.JSON(200, basic.Success("ok"))

}
