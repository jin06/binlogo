package filter

import (
	"context"
	message2 "github.com/jin06/binlogo/app/pipeline/message"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"testing"
)

func TestRun(t *testing.T) {
	pipe := pipeline.Pipeline{
		Name: "test", Filters: []*pipeline.Filter{
			&pipeline.Filter{
				Type: pipeline.FILTER_BLACK,
				Rule: "mysql",
			},
		},
	}

	f, err  := New(WithPipe(&pipe))
	f.InChan =  make(chan *message2.Message, 1)
	f.OutChan=  make(chan *message2.Message, 1)
	if err != nil {
		t.Error(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	err = f.Run(ctx)
	if err  != nil {
		t.Error(err)
	}
	inMsg := &message2.Message{
		Content: &message2.Content{
			Head: &message2.Head{
				Database: "mysql",
				Table:    "user",
			},
		},
	}

	f.InChan <- inMsg

	outMsg := <- f.OutChan
	if outMsg.Filter == false {
		t.Fail()
	}
	cancel()
}


