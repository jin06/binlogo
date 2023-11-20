package http

import (
	"github.com/go-resty/resty/v2"
	message2 "github.com/jin06/binlogo/app/pipeline/message"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
)

// Http send message to http api
type Http struct {
	Http   *pipeline.Http
	Client *resty.Client
}

// New returns a new Http
func New(cfg *pipeline.Http) (h *Http, err error) {
	if cfg.Retries < 0 {
		cfg.Retries = 0
	}
	h = &Http{
		Http: cfg,
	}
	h.Client = resty.New()
	return
}

// Send logic and control
func (h *Http) Send(msg *message2.Message) (ok bool, err error) {
	for i := 0; i <= h.Http.Retries; i++ {
		ok, err = h.doSend(msg)
		if err != nil {
			continue
		}
		if !ok {
			continue
		}
	}
	return
}

type httpResult struct {
	Code int    `json:"code"`
	Data string `json:"data"`
}

func (h *Http) doSend(msg *message2.Message) (ok bool, err error) {
	result := &httpResult{}
	_, err = h.Client.
		R().
		SetHeader("Content-Type", "application/json").
		SetBody(msg.Content).
		SetResult(result).
		Post(h.Http.API)
	if err != nil {
		return
	}
	if result.Code == 2000 {
		ok = true
	}
	return
}

func (h *Http) Close() error {
	return nil
}
