package sender

import "github.com/jin06/binlogo/pipeline/message"

type Sender interface {
	Send(ch *message.Message) (bool, error)
}
