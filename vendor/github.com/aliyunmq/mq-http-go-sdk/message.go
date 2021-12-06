package mq_http_sdk

import (
	"encoding/xml"
	"github.com/gogap/errors"
	"strconv"
	"strings"
)

const (
	StartDeliverTime = "__STARTDELIVERTIME"
	TransCheckImmunityTime = "__TransCheckT"
	Keys = "KEYS"
	SHARDING = "__SHARDINGKEY"
)

type MessageResponse struct {
	XMLName xml.Name `xml:"Message" json:"-"`
	// 错误码
	Code string `xml:"Code,omitempty" json:"code,omitempty"`
	// 错误描述
	Message string `xml:"Message,omitempty" json:"message,omitempty"`
	// 请求ID
	RequestId string `xml:"RequestId,omitempty" json:"request_id,omitempty"`
	// 请求HOST
	HostId string `xml:"HostId,omitempty" json:"host_id,omitempty"`
}

type ErrorResponse struct {
	XMLName xml.Name `xml:"Error" json:"-"`
	// 错误码
	Code string `xml:"Code,omitempty" json:"code,omitempty"`
	// 错误描述
	Message string `xml:"Message,omitempty" json:"message,omitempty"`
	// 请求ID
	RequestId string `xml:"RequestId,omitempty" json:"request_id,omitempty"`
	// 请求HOST
	HostId string `xml:"HostId,omitempty" json:"host_id,omitempty"`
}

type PublishMessageRequest struct {
	XMLName xml.Name `xml:"Message" json:"-"`
	// 消息体
	MessageBody string `xml:"MessageBody" json:"message_body"`
	// 消息标签
	MessageTag string `xml:"MessageTag,omitempty" json:"message_tag,omitempty"`
	// 消息属性
	Properties map[string]string `xml:"-" json:"-"`
	// 消息KEY
	MessageKey string `xml:"-" json:"-"`
	// 定时消息，单位毫秒（ms）,在指定时间戳（当前时间之后）进行投递
	StartDeliverTime int64 `xml:"-" json:"-"`
	// 在消息属性中添加第一次消息回查的最快时间,单位秒,并且表征这是一条事务消息, 10~300s
	TransCheckImmunityTime int `xml:"-" json:"-"`
	// 分区顺序消息中区分不同分区的关键字段，sharding key 于普通消息的 key 是完全不同的概念。全局顺序消息，该字段可以设置为任意非空字符串。
	ShardingKey string `xml:"-" json:"-"`
	// 序列化属性请勿使用
	string `xml:"Properties,omitempty" json:"properties,omitempty"`
}

type PublishMessageResponse struct {
	MessageResponse
	// 消息ID
	MessageId string `xml:"MessageId" json:"message_id"`
	// 消息体MD5
	MessageBodyMD5 string `xml:"MessageBodyMD5" json:"message_body_md5"`
	// 消息句柄只有事务消息才有
	ReceiptHandle string `xml:"ReceiptHandle" json:"receipt_handle"`
}

type ReceiptHandles struct {
	XMLName xml.Name `xml:"ReceiptHandles"`
	// 消息句柄
	ReceiptHandles []string `xml:"ReceiptHandle"`
}

type AckMessageErrorEntry struct {
	XMLName xml.Name `xml:"Error" json:"-"`
	// Ack消息出错的错误码
	ErrorCode string `xml:"ErrorCode" json:"error_code"`
	// Ack消息出错的错误描述
	ErrorMessage string `xml:"ErrorMessage" json:"error_messages"`
	// Ack消息出错的消息句柄
	ReceiptHandle string `xml:"ReceiptHandle,omitempty" json:"receipt_handle"`
}

type AckMessageErrorResponse struct {
	XMLName        xml.Name               `xml:"Errors" json:"-"`
	FailedMessages []AckMessageErrorEntry `xml:"Error" json:"errors"`
}

type ConsumeMessageEntry struct {
	MessageResponse
	// 消息ID
	MessageId string `xml:"MessageId" json:"message_id"`
	// 消息句柄
	ReceiptHandle string `xml:"ReceiptHandle" json:"receipt_handle"`
	// 消息体MD5
	MessageBodyMD5 string `xml:"MessageBodyMD5" json:"message_body_md5"`
	// 消息体
	MessageBody string `xml:"MessageBody" json:"message_body"`
	// 消息发送时间
	PublishTime int64 `xml:"PublishTime" json:"publish_time"`
	// 下次消费消息的时间（如果这次消费的消息没有Ack）
	NextConsumeTime int64 `xml:"NextConsumeTime" json:"next_consume_time"`
	// 第一次消费的时间,此值对于顺序消费没有意义
	FirstConsumeTime int64 `xml:"FirstConsumeTime" json:"first_consume_time"`
	// 消费次数
	ConsumedTimes int64 `xml:"ConsumedTimes" json:"consumed_times"`
	// 消息标签
	MessageTag string `xml:"MessageTag" json:"message_tag"`
	// 消息属性
	Properties map[string]string `xml:"-" json:"-"`
	// 序列化属性请勿使用
	PropInner string `xml:"Properties,omitempty" json:"properties,omitempty"`
	// 消息KEY
	MessageKey string `xml:"-" json:"-"`
	// 定时消息，单位毫秒（ms）,在指定时间戳（当前时间之后）进行投递
	StartDeliverTime int64 `xml:"-" json:"-"`
	// 顺序消息分区Key
	ShardingKey string `xml:"-" json:"-"`
	// 在消息属性中添加第一次消息回查的最快时间,单位秒,并且表征这是一条事务消息
	TransCheckImmunityTime int `xml:"-" json:"-"`
}

type ConsumeMessageResponse struct {
	XMLName  xml.Name              `xml:"Messages" json:"-"`
	Messages []ConsumeMessageEntry `xml:"Message" json:"messages"`
}

func ConstructPubMessage(pubMsgReq *PublishMessageRequest) (err error) {
	if pubMsgReq.Properties == nil {
		pubMsgReq.Properties = make(map[string]string)
	}

	if pubMsgReq.MessageKey != "" {
		pubMsgReq.Properties[Keys] = pubMsgReq.MessageKey
	}

	if pubMsgReq.StartDeliverTime > 0 {
		pubMsgReq.Properties[StartDeliverTime] = strconv.FormatInt(pubMsgReq.StartDeliverTime, 10)
	}

	if pubMsgReq.TransCheckImmunityTime > 0 {
		pubMsgReq.Properties[TransCheckImmunityTime] = strconv.Itoa(pubMsgReq.TransCheckImmunityTime)
	}

	if pubMsgReq.ShardingKey != "" {
		pubMsgReq.Properties[SHARDING] = pubMsgReq.ShardingKey
	}

	if pubMsgReq.Properties == nil {
		return nil
	}

	for k, v := range pubMsgReq.Properties {

		if ContainsSpecialChar(k) || ContainsSpecialChar(v) {
			return ErrMessageProperty.New(errors.Params{"err": k + ":" + v})
		}

		pubMsgReq.string += k + ":" + v + "|"
	}

	return nil
}

func ConstructRecMessage(entries *[]ConsumeMessageEntry) {
	for i := range *entries {
		(*entries)[i].Properties = make(map[string]string)

		if (*entries)[i].PropInner == "" {
			continue
		}
		kvArray := strings.Split((*entries)[i].PropInner, "|")
		for _, kv := range kvArray {
			if kv == "" {
				continue
			}
			kAndV := strings.Split(kv, ":")
			if len(kAndV) != 2 || kAndV[0] == "" || kAndV[1] == "" {
				continue
			}
			(*entries)[i].Properties[kAndV[0]] = kAndV[1]
		}

		(*entries)[i].MessageKey = (*entries)[i].Properties[Keys]
		(*entries)[i].StartDeliverTime, _ = strconv.ParseInt((*entries)[i].Properties[StartDeliverTime], 10, 64)
		(*entries)[i].TransCheckImmunityTime, _ = strconv.Atoi((*entries)[i].Properties[TransCheckImmunityTime])
		(*entries)[i].ShardingKey = (*entries)[i].Properties[SHARDING]
		delete((*entries)[i].Properties, Keys)
		delete((*entries)[i].Properties, StartDeliverTime)
		delete((*entries)[i].Properties, TransCheckImmunityTime)
		delete((*entries)[i].Properties, SHARDING)

		(*entries)[i].PropInner = ""
	}

	return
}

func ContainsSpecialChar(input string) (result bool) {
	return strings.Contains(input, "'") ||
		strings.Contains(input, "\"") ||
		strings.Contains(input, "<") ||
		strings.Contains(input, ">") ||
		strings.Contains(input, "&") ||
		strings.Contains(input, ":") ||
		strings.Contains(input, "|")
}
