package node

import (
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/app/server/console/handler"
	"github.com/jin06/binlogo/app/server/console/module/node"
	"github.com/jin06/binlogo/pkg/node/role"
	"github.com/jin06/binlogo/pkg/store/dao/dao_cluster"
	"github.com/jin06/binlogo/pkg/store/dao/dao_node"
	"github.com/sirupsen/logrus"
	"sort"
	"strconv"
	"strings"
)

func List(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	qSort := c.Query("sort")
	name := c.Query("name")
	ready := c.Query("ready")

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
	resItems := []*node.Item{}
	for _, v := range items {
		if ready != "" {
			if ready == "yes" && v.Status.Ready == false {
				continue
			}
			if ready == "no" && v.Status.Ready == true {
				continue
			}
		}
		if name != "" {
			if !strings.Contains(v.Node.Name, name) {
				continue
			}
		}
		resItems = append(resItems, v)
	}

	sort.Slice(resItems, func(i, j int) bool {
		if qSort == "+id" {
			return resItems[i].Node.CreateTime.Before(resItems[j].Node.CreateTime)
		} else {
			return resItems[j].Node.CreateTime.Before(resItems[i].Node.CreateTime)
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
