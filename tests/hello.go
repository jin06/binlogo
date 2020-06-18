package main

import (
	"fmt"
)

func main()  {
	fmt.Println("Test thing .. ")
	s := "lllllisdjfisowe"
	k := 6
	fmt.Println(reverseLeftWords(s, k))
}

func reverseLeftWords(s string, n int) string {
	fmt.Println(s[0])
	return s;
}