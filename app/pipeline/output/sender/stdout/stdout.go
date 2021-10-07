package stdout

import (
	"fmt"
	message2 "github.com/jin06/binlogo/app/pipeline/message"
	"os"
)

type Stdout struct {
}

func New() (std *Stdout, err error) {
	std = &Stdout{}
	return
}

func (s *Stdout) Send(msg *message2.Message) (ok bool, err error) {
	_, err = fmt.Fprintln(os.Stdout, msg.ToString())
	if err == nil {
		ok = true
	}
	return
}
