package output

import "github.com/jin06/binlogo/pipeline/message"

type Output struct {
	InChan chan *message.Message
}

func (o *Output) Start() (err error) {
	return
}
