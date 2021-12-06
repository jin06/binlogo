package rocketmq

import (
	"github.com/aliyunmq/mq-http-go-sdk"
	message2 "github.com/jin06/binlogo/app/pipeline/message"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
)

// RocketMQ send message to Aliyun's RocketMQ
type RocketMQ struct {
	RocketMQ *pipeline.RocketMQ
	client   mq_http_sdk.MQClient
	producer mq_http_sdk.MQProducer
}

// New returns a new RocketMQ
func New(rm *pipeline.RocketMQ) (r *RocketMQ, err error) {
	r = &RocketMQ{
		RocketMQ: rm,
	}
	r.client = mq_http_sdk.NewAliyunMQClient(r.RocketMQ.Endpoint, r.RocketMQ.AccessKey, r.RocketMQ.SecretKey, "")
	r.producer = r.client.GetProducer(r.RocketMQ.InstanceId, r.RocketMQ.TopicName)
	return
}

// Send logic and control
func (r *RocketMQ) Send(msg *message2.Message) (ok bool, err error) {
	content, _ := msg.JsonContent()
	req := mq_http_sdk.PublishMessageRequest{
		MessageBody: content,
		MessageTag:  "bilogo",
		Properties:  map[string]string{},
		ShardingKey: msg.Content.Head.Database + "_" + msg.Content.Head.Table,
	}
	_, err = r.producer.PublishMessage(req)
	if err == nil {
		ok = true
	}
	return
}
