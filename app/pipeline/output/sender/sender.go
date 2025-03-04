package sender

import (
	"github.com/jin06/binlogo/v2/app/pipeline/message"
)

// Sender interface for sender
type Sender interface {
	Send(msg *message.Message) (bool, error)
	Close() error
}
