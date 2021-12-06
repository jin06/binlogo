package mq_http_sdk

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"net/http"
	neturl "net/url"
	"strconv"
	"sync"
	"time"

	"github.com/gogap/errors"
	"github.com/valyala/fasthttp"
)

const (
	ClientName           = "mq-go-sdk/1.0.3(fasthttp)"
	ClientVersion        = "2015-06-06"
	DefaultTimeout int64 = 35
)

type Method string

var (
	errMapping map[string]errors.ErrCodeTemplate
)

const (
	GET    Method = "GET"
	POST          = "POST"
	DELETE        = "DELETE"
)

type MQClient interface {
	GetProducer(instanceId string, topicName string) MQProducer
	GetTransProducer(instanceId string, topicName string, groupId string) MQTransProducer
	GetConsumer(instanceId string, topicName string, consumer string, messageTag string) MQConsumer
	GetCredential() Credential
	GetEndpoint() string
}

type AliyunMQClient struct {
	timeout      time.Duration
	endpoint     *neturl.URL
	credential   Credential
	client       *fasthttp.Client
	clientLocker sync.Mutex
}

func NewAliyunMQClient(endpoint, accessKeyId, accessKeySecret, securityToken string) MQClient {
	return NewAliyunMQClientWithTimeout(endpoint, accessKeyId, accessKeySecret, securityToken, time.Second*time.Duration(DefaultTimeout))
}

func NewAliyunMQClientWithTimeout(endpoint, accessKeyId, accessKeySecret, securityToken string, timeout time.Duration) MQClient {
	if endpoint == "" {
		panic("mq_go_sdk: endpoint is empty")
	}

	credential := NewMQCredential(accessKeyId, accessKeySecret, securityToken)

	cli := new(AliyunMQClient)
	cli.credential = credential
	cli.timeout = timeout

	var err error
	if cli.endpoint, err = neturl.Parse(endpoint); err != nil {
		panic("err parse endpoint")
	}

	// 2. now init http client
	cli.initFastHttpClient()

	return cli
}

func (p *AliyunMQClient) GetProducer(instanceId string, topicName string) MQProducer {
	if topicName == "" {
		panic("mq_go_sdk: topic name could not be empty")
	}

	ret := new(AliyunMQProducer)
	ret.client = p
	ret.topicName = topicName
	ret.instanceId = instanceId
	ret.decoder = NewAliyunMQDecoder()
	return ret
}

func (p *AliyunMQClient) GetTransProducer(instanceId string, topicName string, groupId string) MQTransProducer {
	if topicName == "" {
		panic("mq_go_sdk: topic name could not be empty")
	}

	if groupId == "" {
		panic("mq_go_sdk: group id could not be empty")
	}

	producer := new(AliyunMQProducer)
	producer.client = p
	producer.topicName = topicName
	producer.instanceId = instanceId
	producer.decoder = NewAliyunMQDecoder()

	ret := new(AliyunMQTransProducer)
	ret.producer = producer
	ret.groupId = groupId
	ret.consumeDecoder = NewAliyunMQDecoder()
	return ret
}

func (p *AliyunMQClient) GetConsumer(instanceId string, topicName string, consumer string, messageTag string) MQConsumer {
	if topicName == "" {
		panic("mq_go_sdk: topic name could not be empty")
	}

	if consumer == "" {
		panic("mq_go_sdk: consumer name could not be empty")
	}

	ret := new(AliyunMQConsumer)
	ret.client = p
	ret.instanceId = instanceId
	ret.topicName = topicName
	ret.consumer = consumer
	if messageTag != "" {
		ret.messageTag = neturl.QueryEscape(messageTag)
	}
	ret.decoder = NewAliyunMQDecoder()

	return ret
}

func (p *AliyunMQClient) GetCredential() Credential {
	return p.credential
}

func (p *AliyunMQClient) GetEndpoint() string {
	return p.endpoint.String()
}

func (p *AliyunMQClient) initFastHttpClient() {
	p.clientLocker.Lock()
	defer p.clientLocker.Unlock()

	p.client = &fasthttp.Client{ReadTimeout: p.timeout, WriteTimeout: p.timeout, Name: ClientName}
}

func (p *AliyunMQClient) authorization(method Method, headers map[string]string, resource string) (authHeader string, err error) {
	if signature, e := p.credential.Signature(method, headers, resource); e != nil {
		return "", e
	} else {
		authHeader = fmt.Sprintf("MQ %s:%s", p.credential.AccessKeyId(), signature)
	}

	return
}

func (p *AliyunMQClient) Send0(method Method, headers map[string]string, message interface{}, resource string) (*fasthttp.Response, error) {
	var xmlContent []byte
	var err error

	if message == nil {
		xmlContent = []byte{}
	} else {
		switch m := message.(type) {
		case []byte:
			{
				xmlContent = m
			}
		default:
			if bXml, e := xml.Marshal(message); e != nil {
				err = ErrMarshalMessageFailed.New(errors.Params{"err": e})
				return nil, err
			} else {
				xmlContent = bXml
			}
		}
	}

	xmlMD5 := md5.Sum(xmlContent)
	strMd5 := fmt.Sprintf("%x", xmlMD5)

	if headers == nil {
		headers = make(map[string]string)
	}

	headers[MQ_VERSION] = ClientVersion
	headers[CONTENT_TYPE] = "text/xml;charset=utf-8"
	headers[CONTENT_MD5] = base64.StdEncoding.EncodeToString([]byte(strMd5))
	headers[DATE] = time.Now().UTC().Format(http.TimeFormat)

	if authHeader, e := p.authorization(method, headers, fmt.Sprintf("/%s", resource)); e != nil {
		err = ErrGeneralAuthHeaderFailed.New(errors.Params{"err": e})
		return nil, err
	} else {
		headers[AUTHORIZATION] = authHeader
	}

	if len(p.credential.SecurityToken()) > 0 {
		headers[SECURITY_TOKEN] = p.credential.SecurityToken()
	}

	var buffer bytes.Buffer
	buffer.WriteString(p.endpoint.String())
	buffer.WriteString("/")
	buffer.WriteString(resource)

	url := buffer.String()

	req := fasthttp.AcquireRequest()

	req.SetRequestURI(url)
	req.Header.SetMethod(string(method))
	req.SetBody(xmlContent)

	for header, value := range headers {
		req.Header.Set(header, value)
	}

	resp := fasthttp.AcquireResponse()

	if err = p.client.Do(req, resp); err != nil {
		err = ErrSendRequestFailed.New(errors.Params{"err": err})
		return nil, err
	}

	return resp, nil
}

func (p *AliyunMQClient) Send(decoder MQDecoder, method Method, headers map[string]string, message interface{}, resource string, v interface{}) (statusCode int, err error) {
	var resp *fasthttp.Response
	if resp, err = p.Send0(method, headers, message, resource); err != nil {
		return
	}

	if resp != nil {
		statusCode = resp.Header.StatusCode()

		if statusCode != fasthttp.StatusCreated &&
			statusCode != fasthttp.StatusOK &&
			statusCode != fasthttp.StatusNoContent {

			// get the response body
			// the body is set in error when decoding xml failed
			bodyBytes := resp.Body()

			var e2 error
			err, e2 = decoder.DecodeError(bodyBytes, resource, string(resp.Header.Peek("x-mq-request-id")))

			if e2 != nil {
				err = ErrUnmarshalErrorResponseFailed.New(errors.Params{"err": e2, "resp": string(bodyBytes)})
				return
			}
			return
		}

		if v != nil {
			buf := bytes.NewReader(resp.Body())
			if e := decoder.Decode(buf, v); e != nil {
				err = ErrUnmarshalResponseFailed.New(errors.Params{"err": e})
				return
			}
			responder := v.(MessageResponder)
			messageRespSlice := responder.Responses()
			for i := range messageRespSlice {
				messageRespSlice[i].Code = strconv.Itoa(resp.StatusCode())
				messageRespSlice[i].RequestId = string(resp.Header.Peek("x-mq-request-id"))
				messageRespSlice[i].HostId = resp.RemoteAddr().String()
			}
		}
	}

	return
}

func ParseError(resp ErrorResponse, resource string) (err error) {
	err = ErrMqServer.New(errors.Params{"resp": resp, "resource": resource})
	return
}
