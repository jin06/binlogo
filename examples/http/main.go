package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/v2/app/pipeline/message"
)

func main() {
	g := gin.Default()

	g.POST("/event", func(c *gin.Context) {
		q := message.Content{}
		if err := c.BindJSON(&q); err != nil {
			panic(123)
		}
		fmt.Println(q)
	})
	fmt.Println(g.Run(":10001"))
}
