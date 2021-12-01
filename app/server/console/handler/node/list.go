package node

import (
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/app/server/console/handler"
	"github.com/jin06/binlogo/app/server/console/module/node"
	"github.com/jin06/binlogo/app/server/console/util"
	"github.com/jin06/binlogo/pkg/node/role"
	"github.com/jin06/binlogo/pkg/store/dao/dao_cluster"
	"github.com/jin06/binlogo/pkg/store/dao/dao_node"
	node2 "github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/sirupsen/logrus"
	"sort"
	"strings"
)

func List(c *gin.Context) {
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
		} else {
			i.Capacity = node2.NewCapacity(v.Name)
		}
		if _, ok := statusMap[v.Name]; ok {
			i.Status = statusMap[v.Name]
		} else {
			i.Status = node2.New(v.Name)
		}
		if i.Node.Name == leaderNode {
			i.Info.Role = role.LEADER
		}
		items = append(items, i)
	}
	resItems := _handleItems(ready, name, qSort, items)

	start, end := util.StartEnd(c)

	if end > len(resItems) {
		end = len(resItems)
	}

	c.JSON(200, handler.Success(map[string]interface{}{
		"items": resItems[start:end],
		"total": len(resItems),
	}))
}

func _handleItems(ready string, name string, qSort string, items []*node.Item) (resItems []*node.Item) {
	resItems = []*node.Item{}
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
	return
}
