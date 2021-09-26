package sender

import "github.com/jin06/binlogo/pipeline/message"

type Sender interface {
	Send(ch chan *message.Message) (bool, error)
}
