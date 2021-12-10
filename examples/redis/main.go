package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	list := "test-redis"
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:16379",
		Password: "",
		DB:       0,
	})
	fmt.Println("Wait...")

	for {
		select {
		case <-time.Tick(duration):
			{
				str, err := rdb.LPop(context.Background(), list).Result()
				if err == redis.Nil {
					fmt.Printf("%v, no message \n", time.Now().Format("2006-01-02 15:04:05.000"))
					duration = Increase(duration)
					continue
				}
				if err != nil {
					rdb.LPush(context.Background(), list, str)
					panic(err)
				}
				duration = Decrease(duration)
				fmt.Println(str)
			}
		}
	}
}

var min = time.Millisecond
var max = time.Second * 10
var duration = min

func Increase(d time.Duration) time.Duration {
	d = d * 2
	if d > max {
		d = max
	}
	return d
}

func Decrease(d time.Duration) time.Duration {
	d = d / 2
	if d < min {
		d = min
	}
	return d
}
