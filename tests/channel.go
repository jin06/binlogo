package main

import "fmt"

func main() {

	c := make(chan string, 2)

	c <- "A"
	c <- "B"
	go func() {
		fmt.Println( <- c)
		fmt.Println( <- c)
	}()
	fmt.Println("write")
	c <- "C"
	c <- "D"
	fmt.Println("end")

}