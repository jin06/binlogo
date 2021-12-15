package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Time{}
	fmt.Println(t.Add(time.Second * 5))
	fmt.Println(t.Add(time.Second * 5).Before(time.Now()))
	t1 := time.Now()
	t2 := t1.Add(time.Second * 5)
	fmt.Println(t1.Before(t2))
	fmt.Println(t)
	fmt.Println(t1)
	fmt.Println(t2)
}
