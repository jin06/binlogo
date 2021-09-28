package stdout

import (
	"fmt"
	"github.com/jin06/binlogo/pipeline/message"
	"os"
)

type Stdout struct {
}

func New() (std *Stdout, err error) {
	std = &Stdout{}
	return
}

func (s *Stdout) Send(msg *message.Message) (ok bool, err error) {
	_, err = fmt.Fprintln(os.Stdout, msg.ToString())
	if err == nil {
		ok = true
	}
	return
}
