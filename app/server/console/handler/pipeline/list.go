package pipeline

import (
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/app/server/console/handler"
	"github.com/jin06/binlogo/app/server/console/module/pipeline"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
	"sort"
	"strconv"
	"strings"
)

func List(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	qSort := c.Query("sort")
	name := c.Query("name")
	status := c.Query("status")

	all, err := dao_pipe.AllPipelines()

	if err != nil {
		c.JSON(200, handler.Fail(err))
		return
	}
	var items []*pipeline.Item
	for _, v := range all {
		if v.IsDelete {
			continue
		}
		items = append(items, &pipeline.Item{Pipeline: v})
	}

	if err = pipeline.CompleteInfoList(items); err != nil {
		c.JSON(200, handler.Fail(err))
		return
	}

	resItems := []*pipeline.Item{}
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

	start := (page - 1) * limit
	if start < 0 {
		start = 0
	}
	end := start + limit
	if end > len(resItems) {
		end = len(resItems)
	}

	c.JSON(200, handler.Success(map[string]interface{}{
		"items": resItems[start:end],
		"total": len(resItems),
	}))
}
