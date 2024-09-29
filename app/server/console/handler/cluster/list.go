package cluster

import (
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/v2/app/server/console/handler"
	"github.com/jin06/binlogo/v2/pkg/store/dao"
	"github.com/jin06/binlogo/v2/pkg/store/dao/dao_cluster"
)

func RegisterList(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	all, err := dao.ALLRegisterNodes()
	if err != nil {
		c.JSON(200, handler.Fail(err))
		return
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].CreateTime.Before(all[j].CreateTime)
	})
	start := (page - 1) * limit
	if start < 0 {
		start = 0
	}
	end := start + limit
	if end > len(all) {
		end = len(all)
	}

	c.JSON(200, handler.Success(map[string]interface{}{
		"items": all[start:end],
		"total": len(all),
	}))

}

func ElectionList(c *gin.Context) {
	all, err := dao_cluster.AllElections()
	if err != nil {
		c.JSON(200, handler.Fail(err))
		return
	}
	sort.Slice(all, func(i, j int) bool {
		valI, ok1 := all[i]["create_revision"].(int64)
		valJ, ok2 := all[j]["create_revision"].(int64)
		if ok1 && ok2 {
			return valI < valJ
		}
		return false
	})
	c.JSON(200, handler.Success(map[string]interface{}{
		"items": all,
		"total": len(all),
	}))
}
