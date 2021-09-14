package input

import "github.com/jin06/binlogo/pipeline/message"

type Controller struct {
	Input Input
}

func (c *Controller) Start() error {
	return c.Input.Start()
}

func (c *Controller) DataLine() chan message.Message {
	return c.Input.DataLine()
}
