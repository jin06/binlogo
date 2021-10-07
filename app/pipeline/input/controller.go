package input

import (
	message2 "github.com/jin06/binlogo/app/pipeline/message"
)

type Controller struct {
	Input Input
}

func (c *Controller) Start() error {
	return c.Input.Run()
}

func (c *Controller) DataLine() chan *message2.Message {
	return c.Input.DataLine()
}
