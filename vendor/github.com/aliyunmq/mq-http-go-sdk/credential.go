package mq_http_sdk

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gogap/errors"
)

const (
	AUTHORIZATION  = "Authorization"
	CONTENT_TYPE   = "Content-Type"
	CONTENT_MD5    = "Content-MD5"
	MQ_VERSION     = "x-mq-version"
	HOST           = "Host"
	DATE           = "Date"
	KEEP_ALIVE     = "Keep-Alive"
	SECURITY_TOKEN = "security-token"
)

type Credential interface {
	Signature(method Method, headers map[string]string, resource string) (signature string, err error)
	AccessKeyId() string
	AccessKeySecret() string
	SecurityToken() string
}

type MQCredential struct {
	accessKeySecret string
	accessKeyId     string
	securityToken   string
}

func NewMQCredential(accessKeyId string, accessKeySecret string, securityToken string) *MQCredential {
	if accessKeyId == "" {
		panic("mq_go_sdk: access key id is empty")
	}
	if accessKeySecret == "" {
		panic("mq_go_sdk: access key secret is empty")
	}
	mqCredential := new(MQCredential)
	mqCredential.accessKeySecret = accessKeySecret
	mqCredential.accessKeyId = accessKeyId
	mqCredential.securityToken = securityToken
	return mqCredential
}

func (p *MQCredential) AccessKeyId() string {
	return p.accessKeyId
}
func (p *MQCredential) AccessKeySecret() string {
	return p.accessKeySecret
}
func (p *MQCredential) SecurityToken() string {
	return p.securityToken
}

func (p *MQCredential) Signature(method Method, headers map[string]string, resource string) (signature string, err error) {
	signItems := []string{}
	signItems = append(signItems, string(method))

	contentMD5 := ""
	contentType := ""
	date := time.Now().UTC().Format(http.TimeFormat)

	if v, exist := headers[CONTENT_MD5]; exist {
		contentMD5 = v
	}

	if v, exist := headers[CONTENT_TYPE]; exist {
		contentType = v
	}

	if v, exist := headers[DATE]; exist {
		date = v
	}

	mqHeaders := []string{}

	for k, v := range headers {
		if strings.HasPrefix(k, "x-mq-") {
			mqHeaders = append(mqHeaders, k+":"+strings.TrimSpace(v))
		}
	}

	sort.Sort(sort.StringSlice(mqHeaders))

	stringToSign := string(method) + "\n" +
		contentMD5 + "\n" +
		contentType + "\n" +
		date + "\n" +
		strings.Join(mqHeaders, "\n") + "\n" +
		resource

	sha1Hash := hmac.New(sha1.New, []byte(p.accessKeySecret))
	if _, e := sha1Hash.Write([]byte(stringToSign)); e != nil {
		err = ErrSignMessageFailed.New(errors.Params{"err": e})
		return
	}

	signature = base64.StdEncoding.EncodeToString(sha1Hash.Sum(nil))

	return
}
