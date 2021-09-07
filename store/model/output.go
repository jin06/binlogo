package model

import "encoding/json"

type Output struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	PipelineId string `json:"pipeline_id"`
}

func (s *Output) Key() (key string) {
	return "pipeline/" + s.PipelineId + "/output/" + s.ID
}

func (s *Output) Val() (val string) {
	b, _ := json.Marshal(s)
	val = string(b)
	return
}

func (s *Output) Unmarshal(val []byte) (err error) {
	err = json.Unmarshal(val, s)
	return
}
