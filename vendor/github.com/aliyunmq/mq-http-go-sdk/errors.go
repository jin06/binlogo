package mq_http_sdk

import (
	"github.com/gogap/errors"
)

const (
	ALIYUN_MQ_ERR_NS = "MQ"

	ALIYUN_MQ_ERR_TEMPSTR = "aliyun_mq response status error,code: {{.resp.Code}}, message: {{.resp.Message}}, path: {{.resource}}, requestId: {{.resp.RequestId}}"
)

var (
	ErrSignMessageFailed       = errors.TN(ALIYUN_MQ_ERR_NS, 1, "sign message failed, {{.err}}")
	ErrMarshalMessageFailed    = errors.TN(ALIYUN_MQ_ERR_NS, 2, "marshal message filed, {{.err}}")
	ErrGeneralAuthHeaderFailed = errors.TN(ALIYUN_MQ_ERR_NS, 3, "general auth header failed, {{.err}}")
	ErrMessageProperty		   = errors.TN(ALIYUN_MQ_ERR_NS, 4, "message property can not contains:\" ' < > & : |, {{.err}}")

	ErrSendRequestFailed = errors.TN(ALIYUN_MQ_ERR_NS, 5, "send request failed, {{.err}}")

	ErrUnmarshalErrorResponseFailed = errors.TN(ALIYUN_MQ_ERR_NS, 7, "unmarshal error response failed, {{.err}}, ResponseBody: {{.resp}}")
	ErrUnmarshalResponseFailed      = errors.TN(ALIYUN_MQ_ERR_NS, 8, "unmarshal response failed, {{.err}}")

	ErrMqServer = errors.TN(ALIYUN_MQ_ERR_NS, 101, ALIYUN_MQ_ERR_TEMPSTR)

	ErrAckMessage = errors.TN(ALIYUN_MQ_ERR_NS, 102, "aliyun_mq ack message error,code: {{.Code}}, message: {{.Message}}, receiptHandle: {{.ReceiptHandle}}, requestId: {{.RequestId}}")
)

type ErrAckItem struct {
	ErrorHandle string
	ErrorMsg    string
	ErrorCode   string
}
