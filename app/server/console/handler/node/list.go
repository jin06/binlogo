package node

import (
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/app/server/console/handler"
	"github.com/jin06/binlogo/app/server/console/module/node"
	"github.com/jin06/binlogo/pkg/store/dao"
)

func List(c *gin.Context) {
	q := struct {
		Name   string `json:"name"`
		Ip     string `json:"ip"`
		Status string `json:"status"`
	}{}

	if err := c.Bind(&q); err != nil {
		c.JSON(200, handler.Fail(err))
		return
	}
	all, err := dao.AllNodes()
	if err != nil {
		c.JSON(200, handler.Fail(err))
		return
	}
	var items []*node.Item
	for _, v := range all {
		items = append(items, &node.Item{Node: v})
	}
	c.JSON(200, handler.Success(map[string]interface{}{
		"items": items,
		"total": len(items),
	}))
}
