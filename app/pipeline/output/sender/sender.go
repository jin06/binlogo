package sender

import (
	"github.com/jin06/binlogo/app/pipeline/message"
)

// Sender interface for sender
type Sender interface {
	Send(ch *message.Message) (bool, error)
	Close() error
}
