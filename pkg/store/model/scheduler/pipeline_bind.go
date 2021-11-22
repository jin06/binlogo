package scheduler

import "encoding/json"

type PipelineBind struct {
	Bindings map[string]string `json:"bindings"`
}

func EmptyPipelineBind() *PipelineBind {
	return &PipelineBind{
		Bindings: map[string]string{},
	}
}

func (s *PipelineBind) Key() (key string) {
	return "scheduler/pipeline_bind"
}

func (s *PipelineBind) Val() (val string) {
	b, _ := json.Marshal(s)
	val = string(b)
	return
}

func (s *PipelineBind) Unmarshal(val []byte) (err error) {
	err = json.Unmarshal(val, s)
	return
}
