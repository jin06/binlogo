package pipeline

import "time"

type Instance struct {
	PipelineName string    `json:"pipeline_name"`
	NodeName     string    `json:"node_name"`
	CreateTime   time.Time `json:"create_time"`
}
