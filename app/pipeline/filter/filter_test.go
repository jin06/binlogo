package filter

import (
	"context"
	"testing"

	"github.com/jin06/binlogo/v2/app/pipeline/message"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
)

func TestRun(t *testing.T) {
	pipe := pipeline.Pipeline{
		Name: "test", Filters: []*pipeline.Filter{
			{
				Type: pipeline.FILTER_BLACK,
				Rule: "mysql",
			},
		},
	}

	f, err := New(WithPipe(&pipe))
	f.InChan = make(chan *message.Message, 1)
	f.OutChan = make(chan *message.Message, 1)
	if err != nil {
		t.Error(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	err = f.Run(ctx)
	if err != nil {
		t.Error(err)
	}
	inMsg := &message.Message{
		Content: message.Content{
			Head: message.Head{
				Database: "mysql",
				Table:    "user",
			},
		},
	}

	f.InChan <- inMsg

	outMsg := <-f.OutChan
	if outMsg.Filter == false {
		t.Fail()
	}
	cancel()
}
