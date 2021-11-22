package stdout

import (
	"encoding/json"
	"fmt"
	message2 "github.com/jin06/binlogo/app/pipeline/message"
	"github.com/jin06/binlogo/pkg/blog"
	"os"
)

type Stdout struct {
}

func New() (std *Stdout, err error) {
	std = &Stdout{}
	return
}

func (s *Stdout) Send(msg *message2.Message) (bool, error) {
	_, err := fmt.Fprintln(os.Stdout, msg.ToString())
	if err != nil {
		blog.Error(err)
	}
	b, err := json.Marshal(msg.Content)
	if err != nil {
		blog.Error(err)
	}
	_, err = fmt.Fprintf(os.Stdout, "Content json: %s \n", string(b))
	if err != nil {
		blog.Error(err)
	}

	return true, nil
}
