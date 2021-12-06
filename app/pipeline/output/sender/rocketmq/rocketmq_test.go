package rocketmq

import (
	message2 "github.com/jin06/binlogo/app/pipeline/message"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"testing"
)

func TestRocketMQ(t *testing.T) {
	rm := &pipeline.RocketMQ{
		Endpoint:   "",
		TopicName:  "",
		InstanceId: "",
		AccessKey:  "",
		SecretKey:  "",
	}
	rocket, err := New(rm)
	if err != nil {
		t.Error(err)
	}
	msg := message2.New()
	msg.Content.Head.Database = "tdb"
	msg.Content.Head.Table = "tbl"
	msg.Content.Data = map[string]string{}
	_, err = rocket.Send(msg)
	if err != nil {
		t.Error(err)
	}
}
