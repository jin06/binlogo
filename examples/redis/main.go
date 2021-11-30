package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

func main() {
	list := "test-redis"
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:16379",
		Password: "",
		DB:       0,
	})
	fmt.Println("Wait...")
	for range time.Tick(time.Second) {
		str, err := rdb.LPop(context.Background(), list).Result()
		if err == redis.Nil {
			fmt.Println("No message")
			continue
		}
		if err != nil {
			rdb.LPush(context.Background(), list, str)
			panic(err)
		}
		fmt.Println(str)
	}
}
