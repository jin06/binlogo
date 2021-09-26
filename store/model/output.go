package model

type Output struct {
	Type       string  `json:"type"`
	PipelineId string  `json:"pipeline_id"`
	Sender     *Sender `json:"sender"`
}
