package mq_http_sdk

import (
	"bytes"
	"encoding/xml"
	"github.com/gogap/errors"
	"io"
)

type MQDecoder interface {
	Decode(reader io.Reader, v interface{}) (err error)
	DecodeError(bodyBytes []byte, resource string, requestId string) (decodedError error, err error)
}

type AliyunMQDecoder struct {
}

type AliyunMQAckMsgDecoder struct {
	v AckMessageErrorResponse
}

func NewAliyunMQDecoder() MQDecoder {
	return &AliyunMQDecoder{}
}

func (p *AliyunMQDecoder) Decode(reader io.Reader, v interface{}) (err error) {
	decoder := xml.NewDecoder(reader)
	err = decoder.Decode(&v)

	return
}

func (p *AliyunMQDecoder) DecodeError(bodyBytes []byte, resource string, requestId string) (decodedError error, err error) {
	bodyReader := bytes.NewReader(bodyBytes)
	errResp := ErrorResponse{}

	decoder := xml.NewDecoder(bodyReader)
	err = decoder.Decode(&errResp)
	if err == nil {
		decodedError = ParseError(errResp, resource)
	}
	return
}

func NewAliyunMQAckMsgDecoder() MQDecoder {
	return &AliyunMQAckMsgDecoder{}
}

func (p *AliyunMQAckMsgDecoder) Decode(reader io.Reader, v interface{}) (err error) {
	decoder := xml.NewDecoder(reader)
	err = decoder.Decode(&v)

	if err == io.EOF {
		err = nil
	}

	return
}

func (p *AliyunMQAckMsgDecoder) DecodeError(bodyBytes []byte, resource string, requestId string) (decodedError error, err error) {
	bodyReader := bytes.NewReader(bodyBytes)

	decoder := xml.NewDecoder(bodyReader)
	err = decoder.Decode(&p.v)
	if err != nil {
		// not receipt handle error
		bodyReader.Seek(0, 0)
		errResp := ErrorResponse{}
		err = decoder.Decode(&errResp)
		if err == nil {
			decodedError = ParseError(errResp, resource)
		}
	} else {
		var errorItems []ErrAckItem
		for _, v := range p.v.FailedMessages {
			errorItems = append(errorItems, ErrAckItem{
				ErrorCode:   v.ErrorCode,
				ErrorHandle: v.ReceiptHandle,
				ErrorMsg:    v.ErrorMessage,
			})
		}
		decodedError = ErrAckMessage.New(errors.Params{
			"Code":          p.v.FailedMessages[0].ErrorCode,
			"Message":       p.v.FailedMessages[0].ErrorMessage,
			"ReceiptHandle": p.v.FailedMessages[0].ReceiptHandle,
			"RequestId":     requestId,
		}).WithContext("Detail", errorItems)
	}
	return
}
