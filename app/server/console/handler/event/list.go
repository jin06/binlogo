package event

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/app/server/console/handler"
	"github.com/jin06/binlogo/pkg/store/dao/dao_event"
)

func List(c *gin.Context) {
	resType := c.Query("res_type")
	resName := c.Query("res_name")
	list, err := dao_event.List(
		resType,
		resName,
		clientv3.WithPrefix(),
		clientv3.WithLimit(20),
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend),
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
