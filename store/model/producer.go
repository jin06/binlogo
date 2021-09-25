package model

import (
	"encoding/json"
	"fmt"
)

type Producer struct {
	ID         string `json:"id"`
	PipelineId string `json:"pipeline_id"`
	OutputId   string `json:"output_id"`
}

func (s *Producer) Key() (key string) {
	return fmt.Sprintf("pipeline/%s/output/%s/producer/%s", s.PipelineId, s.OutputId)
}

func (s *Producer) Val() (val string) {
	b , _ := json.Marshal(s)
	return string(b)
}
