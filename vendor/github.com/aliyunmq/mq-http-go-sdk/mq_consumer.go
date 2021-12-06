package mq_http_sdk

import (
	"fmt"
	"strings"
)

var (
	DefaultNumOfMessages int32 = 16
)

// MQ的消息消费者
type MQConsumer interface {
	// 主题名字
	TopicName() string
	// 实例ID，可空
	InstanceId() string
	// 消费者的名字
	Consumer() string
	// 消费消息过滤的标签
	MessageTag() string
	// 消费消息，如果该条消息没有被 {AckMessage} 确认消费成功，即在NextConsumeTime时会再次消费到该条消息
	ConsumeMessage(respChan chan ConsumeMessageResponse, errChan chan error, numOfMessages int32, waitseconds int64)
	// 顺序消费消息，获取的消息可能是多个分区但是一个分区内消息一定是顺序的，对于一个分区的消息必须要全部ACK成功才能消费下一批消息
	// 否则在NextConsumeTime后会再次消费到相同的消息
	ConsumeMessageOrderly(respChan chan ConsumeMessageResponse, errChan chan error, numOfMessages int32, waitseconds int64)
	// 确认消息消费成功
	AckMessage(receiptHandles []string) (err error)
}

type AliyunMQConsumer struct {
	instanceId string
	topicName  string
	consumer   string
	messageTag string
	client     *AliyunMQClient
	decoder    MQDecoder
}

func (p *AliyunMQConsumer) TopicName() string {
	return p.topicName
}

func (p *AliyunMQConsumer) InstanceId() string {
	return p.instanceId
}

func (p *AliyunMQConsumer) Consumer() string {
	return p.consumer
}

func (p *AliyunMQConsumer) MessageTag() string {
	return p.messageTag
}

func (p *AliyunMQConsumer) ConsumeMessage(respChan chan ConsumeMessageResponse, errChan chan error, numOfMessages int32, waitseconds int64) {
	if numOfMessages <= 0 {
		numOfMessages = DefaultNumOfMessages
	}

	if numOfMessages > 16 {
		numOfMessages = DefaultNumOfMessages
	}

	resourceBuilder := strings.Builder{}

	if p.instanceId != "" {
		resourceBuilder.WriteString(fmt.Sprintf("topics/%s/%s?consumer=%s&ns=%s&numOfMessages=%d",
			p.topicName, "messages", p.consumer, p.instanceId, numOfMessages))
	} else {
		resourceBuilder.WriteString(fmt.Sprintf("topics/%s/%s?consumer=%s&numOfMessages=%d",
			p.topicName, "messages", p.consumer, numOfMessages))
	}

	if waitseconds > 0 && waitseconds <= 30 {
		resourceBuilder.WriteString(fmt.Sprintf("&waitseconds=%d", waitseconds))
	}

	if p.messageTag != "" && len(p.messageTag) > 0 {
		resourceBuilder.WriteString(fmt.Sprintf("&tag=%s", p.messageTag))
	}

	resp := ConsumeMessageResponse{}
	_, err := p.client.Send(p.decoder, GET, nil, nil, resourceBuilder.String(), &resp)
	if err != nil {
		errChan <- err
	} else {
		ConstructRecMessage(&resp.Messages)
		respChan <- resp
	}

	return
}

func (p *AliyunMQConsumer) ConsumeMessageOrderly(respChan chan ConsumeMessageResponse, errChan chan error, numOfMessages int32, waitseconds int64) {
	if numOfMessages <= 0 {
		numOfMessages = DefaultNumOfMessages
	}

	if numOfMessages > 16 {
		numOfMessages = DefaultNumOfMessages
	}

	resourceBuilder := strings.Builder{}

	if p.instanceId != "" {
		resourceBuilder.WriteString(fmt.Sprintf("topics/%s/%s?consumer=%s&trans=order&ns=%s&numOfMessages=%d",
			p.topicName, "messages", p.consumer, p.instanceId, numOfMessages))
	} else {
		resourceBuilder.WriteString(fmt.Sprintf("topics/%s/%s?consumer=%s&trans=order&numOfMessages=%d",
			p.topicName, "messages", p.consumer, numOfMessages))
	}

	if waitseconds > 0 && waitseconds <= 30 {
		resourceBuilder.WriteString(fmt.Sprintf("&waitseconds=%d", waitseconds))
	}

	if p.messageTag != "" && len(p.messageTag) > 0 {
		resourceBuilder.WriteString(fmt.Sprintf("&tag=%s", p.messageTag))
	}

	resp := ConsumeMessageResponse{}
	_, err := p.client.Send(p.decoder, GET, nil, nil, resourceBuilder.String(), &resp)
	if err != nil {
		errChan <- err
	} else {
		ConstructRecMessage(&resp.Messages)
		respChan <- resp
	}

	return
}

func (p *AliyunMQConsumer) AckMessage(receiptHandles []string) (err error) {
	if receiptHandles == nil || len(receiptHandles) == 0 {
		return
	}

	handlers := ReceiptHandles{}

	for _, handler := range receiptHandles {
		handlers.ReceiptHandles = append(handlers.ReceiptHandles, handler)
	}

	resourceBuilder := strings.Builder{}
	if p.instanceId != "" {
		resourceBuilder.WriteString(fmt.Sprintf("topics/%s/%s?consumer=%s&ns=%s", p.topicName, "messages", p.consumer, p.instanceId))
	} else {
		resourceBuilder.WriteString(fmt.Sprintf("topics/%s/%s?consumer=%s", p.topicName, "messages", p.consumer))
	}

	_, err = p.client.Send(NewAliyunMQAckMsgDecoder(), DELETE, nil, handlers, resourceBuilder.String(), nil)

	return
}
