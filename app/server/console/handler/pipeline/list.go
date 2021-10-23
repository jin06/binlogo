package pipeline

import (
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/app/server/console/handler"
	"github.com/jin06/binlogo/pkg/store/dao"
)

func List(c *gin.Context) {
	q := struct {
		Name     string `json:"name"`
		AliasZh  string `json:"alias_zh"`
		SrvType  string `json:"service_type"`
		Describe string `json:"describe"`
	}{}
	if err := c.Bind(&q); err != nil {
		c.JSON(200, handler.Fail(err))
		return
	}
	all, err := dao.AllPipelines()
	if err != nil {
		c.JSON(200, handler.Fail(err))
		return
	}
	var items []*item
	for _, v := range all {
		items = append(items, &item{Pipeline:v})
	}

	c.JSON(200, handler.Success(map[string]interface{}{
		"items": items,
		"total": len(items),
	}))
}
