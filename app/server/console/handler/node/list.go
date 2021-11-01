package node

import (
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/app/server/console/handler"
	"github.com/jin06/binlogo/app/server/console/module/node"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/node/role"
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
	capacityMap, err := dao.CapacityMap()
	if err != nil {
		blog.Error(err)
	}
	leaderNode, err := dao.LeaderNode()
	if err != nil {
		blog.Error(err)
	}
	var items []*node.Item
	for _, v := range all {
		i := &node.Item{
			Node: v,
			Info: &node.Info{
				Role: role.FOLLOWER,
			},
		}
		if _, ok := capacityMap[v.Name]; ok {
			i.Capacity = capacityMap[v.Name]
		}
		if i.Node.Name == leaderNode {
			i.Info.Role = role.LEADER
		}
		items = append(items, i)
	}

	c.JSON(200, handler.Success(map[string]interface{}{
		"items": items,
		"total": len(items),
	}))
}
