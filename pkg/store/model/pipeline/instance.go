package pipeline

import "time"

// Instance store struct
type Instance struct {
	PipelineName string    `json:"pipeline_name"`
	NodeName     string    `json:"node_name"`
	CreateTime   time.Time `json:"create_time"`
}
