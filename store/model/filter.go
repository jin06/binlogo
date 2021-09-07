package model

import "encoding/json"

type Filter struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	PipelineId string `json:"pipeline_id"`
}

func (s *Filter) Key() (key string) {
	return "pipeline/" + s.PipelineId + "/filter/" + s.ID
}

func (s *Filter) Val() (val string) {
	b, _ := json.Marshal(s)
	val = string(b)
	return
}

func (s *Filter) Unmarshal(val []byte) (err error) {
	err = json.Unmarshal(val, s)
	return
}
