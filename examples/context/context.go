package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.TODO())
	go func() {
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			case <-time.Tick(time.Second):
				{
					fmt.Println(1)
				}
			}
		}
	}()
	go func() {
		go func() {
			for {
				select {
				case <-time.Tick(time.Second):
					{
						fmt.Println(3)
					}
				}
			}
		}()
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			case <-time.Tick(time.Second):
				{
					fmt.Println(2)
				}
			}
		}
	}()

	time.Sleep(time.Second * 10)
	cancel()
	select {}
}
