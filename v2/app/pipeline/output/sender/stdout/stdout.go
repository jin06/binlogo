package stdout

import (
	"fmt"

	"github.com/jin06/binlogo/v2/app/pipeline/message"
	"github.com/sirupsen/logrus"
)

// Stdout send message to stdout
// It is mainly used for testing and debugging
type Stdout struct {
}

// New returns a new Stdout
func New() (std *Stdout, err error) {
	std = &Stdout{}
	return
}

// Send logic and control
func (s *Stdout) Send(msg *message.Message) (bool, error) {
	content, err := msg.Json()
	if err != nil {
		logrus.Error(err)
		return true, nil
	}
	fmt.Printf("%s \n", content)
	return true, nil
}

func (s *Stdout) Close() error {
	return nil
}
