package util

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func StartEnd(c *gin.Context) (start int, end int) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	return _StartEnd(page, limit)
}

func _StartEnd(page int, limit int) (start int, end int) {
	if limit <= 0 {
		limit = 10
	}
	start = (page - 1) * limit
	if start < 0 {
		start = 0
	}
	end = start + limit
	return
}
