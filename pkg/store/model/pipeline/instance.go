package pipeline

import (
	"encoding/json"
	"time"
)

// Instance store struct
type Instance struct {
	PipelineName string    `json:"pipeline_name" redis:"pipeline_name"`
	NodeName     string    `json:"node_name" redis:"node_name"`
	CreateTime   time.Time `json:"create_time" redis:"create_time"`
}

func (i *Instance) Val() string {
	b, _ := json.Marshal(i)
	return string(b)
}

func (i *Instance) Unmarshal(val []byte) error {
	return json.Unmarshal(val, i)
}
