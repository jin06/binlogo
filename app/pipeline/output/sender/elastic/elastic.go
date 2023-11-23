package elastic

import (
	"fmt"
	"github.com/jin06/binlogo/app/pipeline/message"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"os"
)

type Elastic struct {
}

func New(cfg *pipeline.Elastic) (e *Elastic, err error) {
	e = &Elastic{}
	return
}

func (e *Elastic) Send(msg *message.Message) (ok bool, err error) {
	ok = true
	_, err = fmt.Fprintf(os.Stdout, "Elastic Send unimplemented")
	return
}

func (e *Elastic) Close() error {
	fmt.Println("Elastic close() unimplemented")
	return nil
}
