package handler

import (
	"github.com/gin-gonic/gin"
	// "github.com/jin06/binlogo/pkg/store/dao/dao_event"
)

// List  event list api
func EventList(c *gin.Context) {
	// resType := c.Query("res_type")
	// resName := c.Query("res_name")
	// list, err := dao.ListEvent(
	// 	resType,
	// 	resName,
	// 	20,
	// 	clientv3.SortByModRevision,
	// 	clientv3.SortDescend,
	// )
	// if err != nil {
	// 	c.JSON(200, basic.Fail(err))
	// 	return
	// }
	// var resList []*model.Event
	// for _, v := range list {
	// 	resList = append(resList, v)
	// }
	// c.JSON(200, basic.Success(map[string]interface{}{
	// 	"items": resList,
	// 	"total": len(resList),
	// }))
}

// ScrollList event scroll api
func EventScrollList(c *gin.Context) {
	// key := c.Query("key")
	// num := c.Query("num")
	// n, err := strconv.Atoi(num)
	// if err != nil {
	// 	n = 20
	// }
	// if n <= 0 {
	// 	n = 20
	// }
	// list, err := dao_event.ScrollList(
	// 	key,
	// 	int64(n),
	// 	clientv3.SortByModRevision,
	// 	clientv3.SortDescend,
	// )

	// if err != nil {
	// 	c.JSON(200, basic.Fail(err))
	// 	return
	// }
	// c.JSON(200, basic.Success(map[string]interface{}{
	// 	"items": list,
	// 	"total": len(list),
	// }))
}
