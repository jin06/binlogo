package model

type Output struct {
	Type       string    `json:"type"`
	PipelineId string    `json:"pipeline_id"`
	Producer   *Producer `json:"producer"`
}
