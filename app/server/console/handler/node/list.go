package node

import (
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/app/server/console/handler"
	"github.com/jin06/binlogo/app/server/console/module/node"
	"github.com/jin06/binlogo/pkg/node/role"
	"github.com/jin06/binlogo/pkg/store/dao/dao_cluster"
	"github.com/jin06/binlogo/pkg/store/dao/dao_node"
	"github.com/sirupsen/logrus"
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
	all, err := dao_node.AllNodes()
	if err != nil {
		c.JSON(200, handler.Fail(err))
		return
	}
	capacityMap, err := dao_node.CapacityMap()
	if err != nil {
		logrus.Error(err)
	}
	statusMap, err := dao_node.StatusMap()
	if err != nil {
		logrus.Error(err)
	}
	leaderNode, err := dao_cluster.LeaderNode()
	if err != nil {
		logrus.Error(err)
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
		if _, ok := statusMap[v.Name]; ok {
			i.Status = statusMap[v.Name]
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
