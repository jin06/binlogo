package scheduler

import "encoding/json"

// PipelineBind pipeline bind
// the node will watch pipeline bind, run pipeline instance if node get the bind,
// stop instance if node lost bind
type PipelineBind struct {
	Bindings map[string]string `json:"bindings"`
}

// EmptyPipelineBind returns a empty pipeline bind
func EmptyPipelineBind() *PipelineBind {
	return &PipelineBind{
		Bindings: map[string]string{},
	}
}

// Key generate etcd prefix of pipeline bind
func (s *PipelineBind) Key() (key string) {
	return "scheduler/pipeline_bind"
}

// Val marshal pipeline bind to json data
func (s *PipelineBind) Val() (val string) {
	b, _ := json.Marshal(s)
	val = string(b)
	return
}

// Unmarshal generate pipeline bind from json data
func (s *PipelineBind) Unmarshal(val []byte) (err error) {
	err = json.Unmarshal(val, s)
	return
}
