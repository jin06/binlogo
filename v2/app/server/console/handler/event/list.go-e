package event

import (
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/app/server/console/handler"
	"github.com/jin06/binlogo/pkg/store/dao/dao_event"
	"github.com/jin06/binlogo/pkg/store/model/event"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strconv"
)

// List  event list api
func List(c *gin.Context) {
	resType := c.Query("res_type")
	resName := c.Query("res_name")
	list, err := dao_event.List(
		resType,
		resName,
		20,
		clientv3.SortByModRevision,
		clientv3.SortDescend,
	)
	if err != nil {
		c.JSON(200, handler.Fail(err))
		return
	}
	var resList []*event.Event
	for _, v := range list {
		resList = append(resList, v)
	}
	c.JSON(200, handler.Success(map[string]interface{}{
		"items": resList,
		"total": len(resList),
	}))
	return
}

// ScrollList event scroll api
func ScrollList(c *gin.Context) {
	key := c.Query("key")
	num := c.Query("num")
	n, err := strconv.Atoi(num)
	if err != nil {
		n = 20
	}
	if n <= 0 {
		n = 20
	}
	list, err := dao_event.ScrollList(
		key,
		int64(n),
		clientv3.SortByModRevision,
		clientv3.SortDescend,
	)

	if err != nil {
		c.JSON(200, handler.Fail(err))
		return
	}
	c.JSON(200, handler.Success(map[string]interface{}{
		"items": list,
		"total": len(list),
	}))
	return
}
