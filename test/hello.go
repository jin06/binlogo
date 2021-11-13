package main

import (
	"fmt"
)

type Model struct {
	Model interface{}
}

func (m *Model)hello()  {
	fmt.Println("model")
}

type Pipeline struct {

}

func (p *Pipeline) hello() {
	fmt.Println("pipeline")
}



func main()  {
	fmt.Println(1<<20)
	fmt.Println("Test thing .. ")
	s := "lllllisdjfisowe"
	k := 6
	fmt.Println(reverseLeftWords(s, k))

	m := Model{}
	m.Model = Pipeline{}
	m.hello()
}

func reverseLeftWords(s string, n int) string {
	fmt.Println(s[0])
	return s
}
