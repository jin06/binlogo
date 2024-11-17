package pipeline

import "time"

// Instance store struct
type Instance struct {
	PipelineName string    `json:"pipeline_name" redis:"pipeline_name"`
	NodeName     string    `json:"node_name" redis:"node_name"`
	CreateTime   time.Time `json:"create_time" redis:"create_time"`
}
