package sender

import (
	message2 "github.com/jin06/binlogo/app/pipeline/message"
)

// Sender interface for sender
type Sender interface {
	Send(ch *message2.Message) (bool, error)
	Close() error
}
