package handler

import (
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/v2/app/server/console/basic"
	"github.com/jin06/binlogo/v2/app/server/console/module/node"
	"github.com/jin06/binlogo/v2/internal/constant"
	"github.com/jin06/binlogo/v2/pkg/store/dao"
	nodeModel "github.com/jin06/binlogo/v2/pkg/store/model/node"
	"github.com/sirupsen/logrus"
)

func NodeList(c *gin.Context) {
	qSort := c.Query("sort")
	name := c.Query("name")
	ready := c.Query("ready")

	all, err := dao.AllNodes(c)
	if err != nil {
		c.JSON(200, basic.Fail(err))
		return
	}
	capacityMap, err := dao.CapacityMap(c)
	if err != nil {
		logrus.Error(err)
	}
	statusMap, err := dao.StatusMap(c)
	if err != nil {
		logrus.Error(err)
	}
	leaderNode, err := dao.LeaderNode(c)
	// leaderNode := "myNode"
	if err != nil {
		logrus.Error(err)
	}
	var items []*node.Item
	for _, v := range all {
		i := &node.Item{
			Node: v,
			Info: &node.Info{
				Role: constant.FOLLOWER,
			},
		}
		if _, ok := capacityMap[v.Name]; ok {
			i.Capacity = capacityMap[v.Name]
		} else {
			i.Capacity = nodeModel.NewCapacity(v.Name)
		}
		if _, ok := statusMap[v.Name]; ok {
			i.Status = statusMap[v.Name]
		} else {
			i.Status = nodeModel.New(v.Name)
		}
		if i.Node.Name == leaderNode {
			i.Info.Role = constant.LEADER
		}
		items = append(items, i)
	}
	resItems := handleItems(ready, name, qSort, items)

	start, end := basic.StartEnd(c)

	if end > len(resItems) {
		end = len(resItems)
	}

	c.JSON(200, basic.Success(map[string]interface{}{
		"items": resItems[start:end],
		"total": len(resItems),
	}))
}

func handleItems(ready string, name string, qSort string, items []*node.Item) (resItems []*node.Item) {
	resItems = []*node.Item{}
	for _, v := range items {
		if ready != "" {
			if ready == "yes" && !v.Status.Ready {
				continue
			}
			if ready == "no" && v.Status.Ready {
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
