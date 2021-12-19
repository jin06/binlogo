package stdout

import (
	"encoding/json"
	"fmt"
	"os"

	message2 "github.com/jin06/binlogo/app/pipeline/message"
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
func (s *Stdout) Send(msg *message2.Message) (bool, error) {
	// _, err := fmt.Fprintln(os.Stdout, msg.ToString())
	// if err != nil {
	// 	logrus.Errorln(err)
	// }
	b, err := json.Marshal(msg.Content)
	if err != nil {
		logrus.Errorln(err)
	}
	_, err = fmt.Fprintf(os.Stdout, "Content json: %s \n", string(b))
	if err != nil {
		logrus.Errorln(err)
	}
	return true, nil
}
