package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main()  {
	for i:=0;i<100;i++ {
		fmt.Println(uint32(rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(100000000)))
	}
}
