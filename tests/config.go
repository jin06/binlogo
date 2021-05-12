package main

import (
	"fmt"
	"github.com/jin06/binlogo/config"
)

func main()  {
	config.InitCfg()
	fmt.Println(config.Cfg)
}