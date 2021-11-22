package pipeline

import (
	"encoding/json"
	"time"
)

type Pipeline struct {
	Name       string    `json:"name"`
	Status     Status    `json:"status"`
	AliasName  string    `json:"aliasName"`
	Mysql      *Mysql    `json:"mysql"`
	Filters    []*Filter `json:"filters"`
	Output     *Output   `json:"output"`
	Replicas   int       `json:"replicas"`
	CreateTime time.Time `json:"create_time"`
	Remark     string    `json:"remark"`
	IsDelete   bool      `json:"is_delete"`
}

type Status string

const (
	STATUS_RUN  Status = "run"
	STATUS_STOP Status = "stop"
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

func (s *Pipeline) ExpectRun() bool {
	if s.IsDelete {
		return false
	}
	if s.Status == STATUS_RUN {
		return true
	}
	if s.Status == STATUS_STOP {
		return false
	}
	return false
}

type OptionPipeline func(p *Pipeline)

func WithPipeStatus(status Status) OptionPipeline {
	return func(p *Pipeline) {
		p.Status = status
	}
}

func WithPipeSafe(uPipe *Pipeline) OptionPipeline {
	return func(p *Pipeline) {
		p.Mysql = uPipe.Mysql
		p.AliasName = uPipe.AliasName
		p.Filters = uPipe.Filters
		p.Output = uPipe.Output
		p.Remark = uPipe.Remark
	}
}

func WithPipeDelete(d bool) OptionPipeline {
	return func(p *Pipeline) {
		p.IsDelete = d
	}
}

func WithPipeMode(mode Mode) OptionPipeline {
	return func(p *Pipeline) {
		p.Mysql.Mode = mode
	}
}

func WithAddFilter(filter *Filter) OptionPipeline {
	return func(p *Pipeline) {
		p.Filters = append(p.Filters, filter)
	}
}

func WithUpdateFilter(index int, filter *Filter) OptionPipeline {
	return func(p *Pipeline) {
		if len(p.Filters) > index {
			p.Filters[index] = filter
		}
	}
}
