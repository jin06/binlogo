package redis

import (
	message2 "github.com/jin06/binlogo/app/pipeline/message"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"testing"
)

func TestRedis(t *testing.T) {
	cfg := pipeline.Redis{
		Addr:     "127.0.0.1:16379",
		UserName: "",
		Password: "",
		DB:       0,
		List:     "go_test_pipeline",
	}
	r, err := New(&cfg)
	if err != nil {
		t.Error(err)
	}
	msg := message2.New()
	ok, err := r.Send(msg)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Fail()
	}
}
