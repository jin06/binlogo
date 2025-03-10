package pipeline

import (
	"encoding/json"
	"time"
)

// Pipeline pipeline's definition
type Pipeline struct {
	Name       string    `json:"name" redis:"name"`
	Status     Status    `json:"status" redis:"status"`
	AliasName  string    `json:"aliasName" redis:"alias_name"`
	Mysql      Mysql     `json:"mysql" redis:"mysql"`
	Filters    Filters   `json:"filters" redis:"filters"`
	Output     Output    `json:"output" redis:"output"`
	Replicas   int       `json:"replicas" redis:"replicas"`
	CreateTime time.Time `json:"create_time" redis:"create_time"`
	Remark     string    `json:"remark" redis:"remark"`
	IsDelete   bool      `json:"is_delete" redis:"is_delete"`
	// If use newest posion to sync mysql replication when get mysql error 1236 (could not find binary log index)
	FixPosNewest bool `json:"fix_pos_newest" redis:"fix_pos_newest"`
}

// NewPipeline returns a new pipeline with default values
func NewPipeline(name string) (pipe *Pipeline) {
	pipe = &Pipeline{
		Name:      name,
		Status:    STATUS_STOP,
		AliasName: name,
		Filters:   Filters{},
		Output: Output{
			Sender: Sender{
				Type: SNEDER_TYPE_STDOUT,
			},
		},
		Replicas:   0,
		CreateTime: time.Now(),
		Remark:     name,
		IsDelete:   false,
	}
	return
}

// Status of Pipeline
type Status = string

const (
	// STATUS_RUN run
	STATUS_RUN Status = "run"
	// STATUS_STOP stop
	STATUS_STOP Status = "stop"
)

// Key generate etcd key
func (s *Pipeline) Key() (key string) {
	return "pipeline/" + s.Name
}

// Val generate json data
func (s *Pipeline) Val() (val string) {
	b, _ := json.Marshal(s)
	val = string(b)
	return
}

// Unmarshal generate from json data
func (s *Pipeline) Unmarshal(val []byte) (err error) {
	err = json.Unmarshal(val, s)
	return
}

// ExpectRun determine whether the pipeline should run
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

// OptionPipeline configure pipeline
type OptionPipeline func(p *Pipeline)

// WithPipeStatus sets pipeline status
func WithPipeStatus(status Status) OptionPipeline {
	return func(p *Pipeline) {
		p.Status = status
	}
}

// WithPipeSafe sets pipeline
func WithPipeSafe(uPipe *Pipeline) OptionPipeline {
	return func(p *Pipeline) {
		p.Mysql = uPipe.Mysql
		p.AliasName = uPipe.AliasName
		p.Filters = uPipe.Filters
		p.Output = uPipe.Output
		p.Remark = uPipe.Remark
		p.FixPosNewest = uPipe.FixPosNewest
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
		p.Filters = append(p.Filters, *filter)
	}
}

func WithUpdateFilter(index int, filter *Filter) OptionPipeline {
	return func(p *Pipeline) {
		if len(p.Filters) > index {
			p.Filters[index] = *filter
		}
	}
}
