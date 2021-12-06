package mq_http_sdk

import (
	"fmt"
	"strings"
)

// MQ的消息生产者
type MQProducer interface {
	// 主题名字
	TopicName() string
	// 实例ID，可空
	InstanceId() string
	// 发送消息
	PublishMessage(message PublishMessageRequest) (resp PublishMessageResponse, err error)
}

// MQ的事务消息生产者
type MQTransProducer interface {
	// 主题名字
	TopicName() string
	// 实例ID，可空
	InstanceId() string
	// GroupId,非空
	GroupId() string
	// 发送消息
	PublishMessage(message PublishMessageRequest) (resp PublishMessageResponse, err error)
	// 消费事务半消息
	ConsumeHalfMessage(respChan chan ConsumeMessageResponse, errChan chan error, numOfMessages int32, waitseconds int64)
	// 提交事务消息
	Commit(receiptHandle string) (err error)
	// 回滚事务消息
	Rollback(receiptHandle string) (err error)
}

type AliyunMQProducer struct {
	topicName  string
	instanceId string
	client     *AliyunMQClient
	decoder    MQDecoder
}

func (p *AliyunMQProducer) TopicName() string {
	return p.topicName
}

func (p *AliyunMQProducer) InstanceId() string {
	return p.instanceId
}

func (p *AliyunMQProducer) PublishMessage(message PublishMessageRequest) (resp PublishMessageResponse, err error) {
	resourceBuilder := strings.Builder{}
	if p.instanceId != "" {
		resourceBuilder.WriteString(fmt.Sprintf("topics/%s/%s?ns=%s", p.topicName, "messages", p.instanceId))
	} else {
		resourceBuilder.WriteString(fmt.Sprintf("topics/%s/%s", p.topicName, "messages"))
	}
	err = ConstructPubMessage(&message)
	if err != nil {
		return
	}
	_, err = p.client.Send(p.decoder, POST, nil, message, resourceBuilder.String(), &resp)
	return
}

type AliyunMQTransProducer struct {
	groupId        string
	producer       *AliyunMQProducer
	consumeDecoder MQDecoder
}

func (p *AliyunMQTransProducer) GroupId() string {
	return p.groupId
}

func (p *AliyunMQTransProducer) TopicName() string {
	return p.producer.topicName
}

func (p *AliyunMQTransProducer) InstanceId() string {
	return p.producer.instanceId
}

func (p *AliyunMQTransProducer) PublishMessage(message PublishMessageRequest) (resp PublishMessageResponse, err error) {
	resp, err = p.producer.PublishMessage(message)
	return
}

func (p *AliyunMQTransProducer) ConsumeHalfMessage(respChan chan ConsumeMessageResponse, errChan chan error, numOfMessages int32, waitseconds int64) {
	if numOfMessages <= 0 {
		numOfMessages = DefaultNumOfMessages
	}

	if numOfMessages > 16 {
		numOfMessages = DefaultNumOfMessages
	}

	resourceBuilder := strings.Builder{}

	if p.InstanceId() != "" {
		resourceBuilder.WriteString(fmt.Sprintf("topics/%s/%s?consumer=%s&ns=%s&numOfMessages=%d&trans=pop",
			p.TopicName(), "messages", p.groupId, p.InstanceId(), numOfMessages))
	} else {
		resourceBuilder.WriteString(fmt.Sprintf("topics/%s/%s?consumer=%s&numOfMessages=%d&trans=pop",
			p.TopicName(), "messages", p.groupId, numOfMessages))
	}

	if waitseconds > 0 && waitseconds <= 30 {
		resourceBuilder.WriteString(fmt.Sprintf("&waitseconds=%d", waitseconds))
	}

	resp := ConsumeMessageResponse{}
	_, err := p.producer.client.Send(p.consumeDecoder, GET, nil, nil, resourceBuilder.String(), &resp)
	if err != nil {
		errChan <- err
	} else {
		ConstructRecMessage(&resp.Messages)
		respChan <- resp
	}

	return
}

func (p *AliyunMQTransProducer) Commit(receiptHandle string) (err error) {
	err = CommitOrRollback(p, receiptHandle, "commit")
	return
}

func (p *AliyunMQTransProducer) Rollback(receiptHandle string) (err error) {
	err = CommitOrRollback(p, receiptHandle, "rollback")
	return
}

func CommitOrRollback(p *AliyunMQTransProducer, receiptHandle string, trans string) (err error) {
	if receiptHandle == "" {
		return
	}

	handlers := ReceiptHandles{}
	handlers.ReceiptHandles = append(handlers.ReceiptHandles, receiptHandle)

	resourceBuilder := strings.Builder{}
	if p.InstanceId() != "" {
		resourceBuilder.WriteString(fmt.Sprintf("topics/%s/%s?consumer=%s&ns=%s&trans=%s",
			p.TopicName(), "messages", p.groupId, p.InstanceId(), trans))
	} else {
		resourceBuilder.WriteString(fmt.Sprintf("topics/%s/%s?consumer=%s&trans=%s",
			p.TopicName(), "messages", p.InstanceId(), trans))
	}

	_, err = p.producer.client.Send(NewAliyunMQAckMsgDecoder(), DELETE, nil, handlers, resourceBuilder.String(), nil)

	return
}

