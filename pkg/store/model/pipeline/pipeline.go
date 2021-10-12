package pipeline

import (
	"encoding/json"
	"github.com/jin06/binlogo/pkg/store/model"
)

type Pipeline struct {
	Name      string          `json:"name"`
	Status    Status          `json:"status"`
	AliasName string          `json:"aliasName"`
	Mysql     *model.Mysql    `json:"mysql"`
	Filters   []*model.Filter `json:"filters"`
	Output    *model.Output   `json:"output"`
	Replicas  int             `json:"replicas"`
}

type Status string

const (
	STATUS_RUN   Status = "run"
	STATUS_STOP  Status = "stop"
)

func (s *Pipeline) Key() (key string) {
	return "pipeline/" + s.Name
}

func (s *Pipeline) Val() (val string) {
	b, _ := json.Marshal(s)
	val = string(b)
	return
}

func (s *Pipeline) Unmarshal(val []byte) (err error) {
	err = json.Unmarshal(val, s)
	return
}

